package database

// SetUsername allows to set a new username for a specified uid.
// Request will fail if username already exists in database
func (db *appdbimpl) SetUsername(userid uint64, newUsername string) error {
	_, err := db.c.Exec("UPDATE user SET username = ? WHERE uid = ?", newUsername, userid)
	return err
}
