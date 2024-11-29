package database

/* MESSAGE */
func (db *appdbimpl) CheckIfMessageExistsByMessageId(messageId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Message WHERE messageId = ?)", messageId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CreateMessage(timestamp string, senderId int, conversationId int, status string, mexType string, content string) error {
	_, err := db.c.Exec("INSERT INTO Message (timestamp, senderId, conversationId, status, type, content, repliedTo, forwardedFrom) VALUES (?, ?, ?, ?, ?, ?, NULL, NULL)", timestamp, senderId, conversationId, status, mexType, content)
	return err
}

func (db *appdbimpl) ReplyToAMessage(timestamp string, senderId int, conversationId int, status string, mexType string, content string, repliedTo int) error {
	_, err := db.c.Exec("INSERT INTO Message (timestamp, senderId, conversationId, status, type, content, repliedTo, forwardedFrom) VALUES (?, ?, ?, ?, ?, ?, ?, NULL)", timestamp, senderId, conversationId, status, mexType, content, repliedTo)
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
	_, err := db.c.Exec("INSERT INTO Message (timestamp, senderId, conversationId, status, type, content, repliedTo, forwardedFrom) VALUES (?, ?, ?, ?, ?, ?, NULL, ?)", message.Timestamp, message.Sender, message.Conversation, message.Status, message.Type, message.Body, messageId)
	return message, err
}

func (db *appdbimpl) UpdateMessageStatus(messageId int, status string) error {
	_, err := db.c.Exec("UPDATE Message SET status = ? WHERE messageId = ?", status, messageId)
	return err
}

func (db *appdbimpl) AddReceiver(userId int, messageId int) error {
	_, err := db.c.Exec("INSERT INTO MessageReceivers(userId, messageId) VALUES (?, ?)", userId, messageId)
	return err
}

func (db *appdbimpl) AddReader(userId int, messageId int) error {
	_, err := db.c.Exec("INSERT INTO MessageReaders(userId, messageId) VALUES (?, ?)", userId, messageId)
	return err
}
