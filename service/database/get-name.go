package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT description FROM GroupConversation WHERE groupId=1").Scan(&name)
	return name, err
}

/* USERS OPERATIONS */
func (db *appdbimpl) CheckIfUserExistsByUsername(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CheckIfUserExistsByUserId(userId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE userId = ?)", userId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) GetUserId(username string) (int, error) {
	var userId int
	err := db.c.QueryRow("SELECT userId FROM User WHERE username = ?", username).Scan(&userId)
	return userId, err
}

func (db *appdbimpl) CreateUser(username string) error {
	_, err := db.c.Exec("INSERT INTO User (username, photoUrl) VALUES (?, ?)", username, "images/default_image.png")
	return err
}

type User struct {
	Username string `json:"username"`
	PhotoUrl string `json:"PhotoUrl"`
}

func (db *appdbimpl) GetUserById(userId int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT username, photoUrl FROM User where userId = ?", userId).Scan(&user.Username, &user.PhotoUrl)
	return user, err
}

func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT username, photoUrl FROM User where username = ?", username).Scan(&user.Username, &user.PhotoUrl)
	return user, err
}

func (db *appdbimpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT username, photoUrl FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username, &user.PhotoUrl); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *appdbimpl) UpdateUsername(username string, userId int) error {
	_, err := db.c.Exec("UPDATE User SET username = ? WHERE userId = ?", username, userId)
	return err
}

func (db *appdbimpl) UpdatePhotoUrl(url string, userId int) error {
	_, err := db.c.Exec("UPDATE User SET photoUrl = ? WHERE userId = ?", url, userId)
	return err
}

/* CONVERSATIONS OPERATIONS */
type ConversationId struct {
	ResourceId int `json:"resourceId"`
}

type Conversation struct {
	Name     string    `json:"conversationName"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Timestamp     string `json:"timestamp"`
	Sender        int    `json:"sender"`
	Conversation  int    `json:"conversation"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Body          string `json:"body"`
	RepliedTo     int    `json:"repliedTo"`
	ForwardedFrom int    `json:"forwardedFrom"`
}

type Comment struct {
	Emoji string `json:"emoji"`
}

func (db *appdbimpl) GetPrivateConversationsByUserId(userId int) ([]ConversationId, error) {
	rows, err := db.c.Query("SELECT conversationId FROM PrivateConversation WHERE userId_1 = ? or userId_2 = ?", userId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []ConversationId
	for rows.Next() {
		var conversation ConversationId
		if err := rows.Scan(&conversation.ResourceId); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	return conversations, nil
}

func (db *appdbimpl) GetGroupConversationsByUserId(userId int) ([]ConversationId, error) {
	rows, err := db.c.Query("SELECT conversationId FROM GroupConversation gc JOIN UserGroup ug ON gc.groupId = ug.groupId WHERE userId = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []ConversationId
	for rows.Next() {
		var conversation ConversationId
		if err := rows.Scan(&conversation.ResourceId); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	return conversations, nil
}

func (db *appdbimpl) CheckIfConversationExistsByConversationId(conversationId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Conversation WHERE conversationId = ?)", conversationId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) GetConversationByConversationId(conversationId int, userId int) (Conversation, error) {
	var conversation Conversation
	var messages []Message
	rows, err := db.c.Query("SELECT timestamp, senderId, status, type, content, repliedTo, forwardedFrom FROM Message WHERE conversationId = ?", conversationId)
	for rows.Next() {
		var message Message
		rows.Scan(&message.Timestamp, &message.Sender, &message.Status, &message.Type, &message.Body, &message.RepliedTo, &message.ForwardedFrom)
		messages = append(messages, message)
	}
	rows.Close()
	conversation.Messages = messages
	var name string
	_ = db.c.QueryRow("SELECT name FROM GroupConversation WHERE conversationId = ?", conversationId).Scan(&name)
	_ = db.c.QueryRow("SELECT u.username FROM PrivateConversation pc JOIN User u ON pc.userId_1 = u.userId OR pc.userId_2 = u.userId WHERE pc.conversationId = ? AND (u.userId != ?)", conversationId, userId).Scan(&name)
	conversation.Name = name
	return conversation, err
}

func (db *appdbimpl) CreatePrivateConversation(userId_1 int, userId_2 int) (int, error) {
	res, _ := db.c.Exec("INSERT INTO Conversation DEFAULT VALUES")
	conversationId, _ := res.LastInsertId()
	_, err := db.c.Exec("INSERT INTO PrivateConversation (conversationId, userId_1, userId_2) VALUES (?, ?, ?)", conversationId, userId_1, userId_2)
	return int(conversationId), err
}

/* MESSAGE */
func (db *appdbimpl) CheckIfMessageExistsByMessageId(messageId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Message WHERE messageId = ?)", messageId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CreateMessage(timestamp string, senderId int, conversationId int, status string, mexType string, content string, repliedTo int, forwardedFrom int) error {
	_, err := db.c.Exec("INSERT INTO Message (timestamp, senderId, conversationId, status, type, content, repliedTo, forwardedFrom) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", timestamp, senderId, conversationId, status, mexType, content, repliedTo, forwardedFrom)
	return err
}

func (db *appdbimpl) DeleteMessage(messageId int) error {
	_, err := db.c.Exec("DELETE FROM Message WHERE messageId = ?", messageId)
	return err
}

func (db *appdbimpl) ForwardMessage(messageId int, senderId int, conversationId int, timestamp string) (Message, error) {
	var message Message
	_ = db.c.QueryRow("SELECT status, type, content FROM Message WHERE messageId = ?", messageId).Scan(&message.Status, &message.Type, &message.Body)
	message.Timestamp = timestamp
	message.Sender = senderId
	message.Conversation = conversationId
	message.ForwardedFrom = messageId
	_, err := db.c.Exec("INSERT INTO Message (timestamp, senderId, conversationId, status, type, content, repliedTo, forwardedFrom) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", message.Timestamp, message.Sender, message.Conversation, message.Status, message.Type, message.Body, 0, messageId)
	return message, err
}
