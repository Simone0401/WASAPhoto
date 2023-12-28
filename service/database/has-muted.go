package database

// HasMuted checks if the specified user uid already muted muteduid user
// Request will fail if specified uid user already muted the muteduid in database
func (db *appdbimpl) HasMuted(userid uint64, muteduid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM ban WHERE uid = ? AND buid = ?", userid, muteduid).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, err
	}

	return true, err
}
