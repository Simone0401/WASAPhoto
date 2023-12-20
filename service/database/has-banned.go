package database

// HasBanned checks if the specified user uid has already banned fuid user
// Request will fail if specified fuid user has already banned the uid in database
func (db *appdbimpl) HasBanned(userid uint64, banneduid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM ban WHERE uid = ? AND buid = ?", userid, banneduid).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, err
	}

	return true, err
}
