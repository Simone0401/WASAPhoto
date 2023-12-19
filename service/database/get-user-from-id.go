package database

// GetUserByID allows to get a User database struct passing an uid.
// Request will fail if uid doesn't exist
func (db *appdbimpl) GetUserByID(uid uint64) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT * FROM user WHERE uid = ?", uid).Scan(&user.Userid, &user.Username)
	return user, err
}
