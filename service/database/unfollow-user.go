package database

// UnfollowUser allows specified uid user to unfollow fuid user
// Request will fail if specified uid user has already unfollowed the fuid in database
// Success request will return true
func (db *appdbimpl) UnfollowUser(userid uint64, followuid uint64) (bool, error) {
	_, err := db.c.Exec("DELETE FROM follow WHERE uid = ? AND fuid = ?", userid, followuid)
	if err != nil {
		return false, err
	}
	return true, err
}
