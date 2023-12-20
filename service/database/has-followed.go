package database

// HasFollowed checks if the specified user uid already follow fuid user
// Request will fail if specified uid user already follow the fuid in database
func (db *appdbimpl) HasFollowed(userid uint64, followuid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM follow WHERE uid = ? AND fuid = ?", userid, followuid).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, err
	}

	return true, err
}
