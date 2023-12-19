package database

// CreateUser allows to create a user from a username.
// Request will fail if username already exists.
func (db *appdbimpl) CreateUser(username string) (User, error) {
	var user User
	var maxId uint64

	// RECOVER MAX(ID)
	err := db.c.QueryRow("SELECT MAX(uid) FROM user").Scan(&maxId)
	if err != nil {
		return User{0, ""}, err
	}

	setId := maxId + 1

	result, err := db.c.Exec("INSERT INTO user (uid, username) VALUES (?, ?)", setId, username)

	if err != nil {
		return User{0, ""}, err
	}

	check, err := result.RowsAffected()

	if err != nil {
		return User{0, ""}, err
	}

	if check > 0 {
		user.Username = username
		user.Userid = setId
	}

	return user, err
}
