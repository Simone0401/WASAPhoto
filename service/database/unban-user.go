package database

// UnbanUser allows specified uid user to unban muteduid user
// Request will fail if specified uid user hasn't already muted the muteduid in database
// Success request will return true
func (db *appdbimpl) UnbanUser(userid uint64, muteduid uint64) (bool, error) {
	_, err := db.c.Exec("DELETE FROM ban WHERE uid = ? AND buid = ?", userid, muteduid)
	if err != nil {
		return false, err
	}
	return true, err
}
