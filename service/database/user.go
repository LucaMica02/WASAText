package database

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

func (db *appdbimpl) GetUserById(userId int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT userId, username, photoUrl FROM User where userId = ?", userId).Scan(&user.ResourceId, &user.Username, &user.PhotoUrl)
	return user, err
}

func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT username, photoUrl FROM User where username = ?", username).Scan(&user.Username, &user.PhotoUrl)
	return user, err
}

func (db *appdbimpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT userId, username, photoUrl FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ResourceId, &user.Username, &user.PhotoUrl); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	// check rows.Err
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *appdbimpl) UpdateUsername(username string, userId int) error {
	_, err := db.c.Exec("UPDATE User SET username = ? WHERE userId = ?", username, userId)
	return err
}

func (db *appdbimpl) UpdateUserPhotoUrl(url string, userId int) error {
	_, err := db.c.Exec("UPDATE User SET photoUrl = ? WHERE userId = ?", url, userId)
	return err
}
