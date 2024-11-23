package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT description FROM GroupConversation WHERE groupId=1").Scan(&name)
	return name, err
}

func (db *appdbimpl) CheckIfUserExists(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) GetUserId(username string) (int, error) {
	var userId int
	err := db.c.QueryRow("SELECT userId FROM User WHERE username = ?", username).Scan(&userId)
	return userId, err
}

func (db *appdbimpl) CreateUser(username string) error {
	_, err := db.c.Exec("INSERT INTO User (username, photoUrl) VALUES (?, ?)", username, "http://default.jpg")
	return err
}
