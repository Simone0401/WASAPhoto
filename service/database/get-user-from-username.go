package database

// GetUserByUsername allows to get a User database struct passing a username.
// Request will fail if username doesn't exist
func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT * FROM user WHERE username = ?", username).Scan(&user.Userid, &user.Username)
	return user, err
}
