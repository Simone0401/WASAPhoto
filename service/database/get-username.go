package database

// GetUsername allows getting the string username value for a specific uid
func (db *appdbimpl) GetUsername(uid uint64) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM user WHERE uid=?`, uid).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, err
}
