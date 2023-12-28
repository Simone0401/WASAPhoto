package database

// BanUser allows specified uid user to ban muteduid user
// Request will fail if specified uid user has already muted the muteduid in database
// Success request will return true
func (db *appdbimpl) BanUser(userid uint64, muteduid uint64) (bool, error) {
	_, err := db.c.Exec("INSERT INTO ban (uid, buid) VALUES (?, ?)", userid, muteduid)
	if err != nil {
		return false, err
	}
	return true, err
}
