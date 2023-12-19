package database

// CheckExistsByUsername checks if a username is already present in database.
// Request will fail if username doesn't exist
func (db *appdbimpl) CheckExistsByUsername(username string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&count)
	return count > 0, err
}
