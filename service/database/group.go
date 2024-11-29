package database

/* GROUP */
func (db *appdbimpl) CheckIfGroupExistsByGroupId(groupId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM GroupConversation WHERE groupId = ?)", groupId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CheckIfUserIsPartecipant(userId int, groupId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM UserGroup WHERE groupId = ? AND userId = ?)", groupId, userId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) CreateGroupConversation(name string, description string, photoUrl string) (int, error) {
	res, _ := db.c.Exec("INSERT INTO Conversation DEFAULT VALUES")
	conversationId, _ := res.LastInsertId()
	res, _ = db.c.Exec("INSERT INTO GroupConversation (name, description, photoUrl, conversationId) VALUES (?, ?, ?, ?)", name, description, photoUrl, conversationId)
	groupId, err := res.LastInsertId()
	return int(groupId), err
}

func (db *appdbimpl) AddUserToGroup(userId int, groupId int) error {
	_, err := db.c.Exec("INSERT INTO UserGroup(userId, groupId) VALUES (?, ?)", userId, groupId)
	return err
}

func (db *appdbimpl) LeaveGroup(userId int, groupId int) error {
	_, err := db.c.Exec("DELETE FROM UserGroup WHERE userId = ? AND groupId = ?", userId, groupId)
	return err
}

func (db *appdbimpl) UpdateGroupName(groupId int, name string) error {
	_, err := db.c.Exec("UPDATE GroupConversation SET name = ? WHERE groupId = ?", name, groupId)
	return err
}

func (db *appdbimpl) UpdateGroupDescription(groupId int, description string) error {
	_, err := db.c.Exec("UPDATE GroupConversation SET description = ? WHERE groupId = ?", description, groupId)
	return err
}

func (db *appdbimpl) UpdateGroupPhotoUrl(url string, groupId int) error {
	_, err := db.c.Exec("UPDATE GroupConversation SET photoUrl = ? WHERE groupId = ?", url, groupId)
	return err
}
