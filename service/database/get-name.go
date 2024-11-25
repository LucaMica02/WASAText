package database

import "time"

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
	Timestamp time.Time `json:"timestamp"`
	Sender    Username  `json:"sender"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Body      string    `json:"body"`
}

type Comment struct {
	Emoji string `json:"emoji"`
}

type Username struct {
	Username string `json:"username"`
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

func (db *appdbimpl) GetConversationByConversationId(conversationId int) (Conversation, error) {
	var conversation Conversation
	var messages []Message
	rows, err := db.c.Query("SELECT timestamp, senderId, status, type, content FROM Message WHERE conversationId = ?", conversationId)
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.Timestamp, &message.Sender, &message.Status, &message.Type, &message.Body)
		messages = append(messages, message)
	}
	rows.Close()
	conversation.Messages = messages
	var name string
	_ = db.c.QueryRow("SELECT name FROM GroupConversation WHERE conversationId = ?", conversationId).Scan(&name)
	conversation.Name = name
	return conversation, err
}
