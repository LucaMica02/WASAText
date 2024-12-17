package database

/* CONVERSATIONS OPERATIONS */

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
	if err != nil {
		return conversation, err
	}
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.Timestamp, &message.Sender, &message.Status, &message.Type, &message.Body, &message.RepliedTo, &message.ForwardedFrom)
		if err != nil {
			return conversation, err
		}
		messages = append(messages, message)
	}
	rows.Close()
	conversation.Messages = messages
	var name string
	err = db.c.QueryRow("SELECT name FROM GroupConversation WHERE conversationId = ?", conversationId).Scan(&name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return conversation, err
	}
	err = db.c.QueryRow("SELECT u.username FROM PrivateConversation pc JOIN User u ON pc.userId_1 = u.userId OR pc.userId_2 = u.userId WHERE pc.conversationId = ? AND (u.userId != ?)", conversationId, userId).Scan(&name)
	conversation.Name = name
	return conversation, err
}

func (db *appdbimpl) CreatePrivateConversation(userId_1 int, userId_2 int) (int, error) {
	res, err := db.c.Exec("INSERT INTO Conversation DEFAULT VALUES") // HERE
	if err != nil {
		return -1, err
	}
	conversationId, _ := res.LastInsertId()
	_, err = db.c.Exec("INSERT INTO PrivateConversation (conversationId, userId_1, userId_2) VALUES (?, ?, ?)", conversationId, userId_1, userId_2)
	return int(conversationId), err
}
