package database

// CheckExistsByUID checks if a user id is already present in database.
// Request will fail if uid doesn't exist
func (db *appdbimpl) CheckExistsByUID(uid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM user WHERE uid = ?", uid).Scan(&count)
	return count > 0, err
}
