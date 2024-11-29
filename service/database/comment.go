package database

/* COMMENT */
func (db *appdbimpl) CheckIfCommentExistsByCommentId(commentId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Comment WHERE commentId = ?)", commentId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CheckIfCommentExists(senderId int, messageId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Comment WHERE senderId = ? AND messageId = ?)", senderId, messageId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CheckIfIsUserComment(senderId int, commentId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Comment WHERE senderId = ? AND commentId = ?)", senderId, commentId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) UpdateComment(senderId int, messageId int, reaction string, timestamp string) error {
	_, err := db.c.Exec("UPDATE Comment SET reaction = ?, timestamp = ? WHERE senderId = ? AND messageId = ?", reaction, timestamp, senderId, messageId)
	return err
}

func (db *appdbimpl) AddComment(timestamp string, senderId int, messageId int, reaction string) error {
	_, err := db.c.Exec("INSERT INTO Comment(timestamp, senderId, messageId, reaction) VALUES (?, ?, ?, ?)", timestamp, senderId, messageId, reaction)
	return err
}

func (db *appdbimpl) DeleteComment(commentId int) error {
	_, err := db.c.Exec("DELETE FROM Comment WHERE commentId = ?", commentId)
	return err
}
