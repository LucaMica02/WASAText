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
	// check rows.Err
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (db *appdbimpl) GetGroupConversationsByUserId(userId int) ([]ConversationId, error) {
	rows, err := db.c.Query("SELECT gc.groupId FROM GroupConversation gc JOIN UserGroup ug ON gc.groupId = ug.groupId WHERE userId = ?", userId)
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
	// check rows.Err
	if err := rows.Err(); err != nil {
		return nil, err
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
	rows, err := db.c.Query("SELECT m.messageId, m.timestamp, m.senderId, m.status, m.type, m.content, m.repliedTo, m.forwardedFrom, count(c.messageId) FROM Message m FULL JOIN Comment c ON m.messageId = c.messageId WHERE conversationId = ? group by m.messageId", conversationId)
	if err != nil {
		return conversation, err
	}
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ResourceId, &message.Timestamp, &message.Sender, &message.Status, &message.Type, &message.Body, &message.RepliedTo, &message.ForwardedFrom, &message.Comments)
		if err != nil {
			return conversation, err
		}
		messages = append(messages, message)
	}
	rows.Close()
	// check rows.Err
	if err := rows.Err(); err != nil {
		return conversation, err
	}
	conversation.Messages = messages
	var name, description, photoUrl string
	var isPrivate bool
	_ = db.c.QueryRow("SELECT name FROM GroupConversation WHERE groupId = ?", conversationId).Scan(&name)
	_ = db.c.QueryRow("SELECT description FROM GroupConversation WHERE groupId = ?", conversationId).Scan(&description)
	err = db.c.QueryRow("SELECT photoUrl FROM GroupConversation WHERE groupId = ?", conversationId).Scan(&photoUrl)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return conversation, err
	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		isPrivate = true
	} else {
		isPrivate = false
	}
	err = db.c.QueryRow("SELECT u.username FROM PrivateConversation pc JOIN User u ON pc.userId_1 = u.userId OR pc.userId_2 = u.userId WHERE pc.conversationId = ? AND (u.userId != ?)", conversationId, userId).Scan(&name)
	conversation.Name = name
	conversation.Description = description
	conversation.IsPrivate = isPrivate
	conversation.PhotoUrl = photoUrl
	return conversation, err
}

func (db *appdbimpl) CreatePrivateConversation(userId_1 int, userId_2 int) (int, error) {
	res, err := db.c.Exec("INSERT INTO Conversation DEFAULT VALUES") 
	if err != nil {
		return -1, err
	}
	conversationId, _ := res.LastInsertId()
	_, err = db.c.Exec("INSERT INTO PrivateConversation (conversationId, userId_1, userId_2) VALUES (?, ?, ?)", conversationId, userId_1, userId_2)
	return int(conversationId), err
}
