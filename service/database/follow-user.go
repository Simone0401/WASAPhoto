package database

// FollowUser allows specified uid user to follow fuid user
// Request will fail if specified uid user has already followed the fuid in database
// Success request will return true
func (db *appdbimpl) FollowUser(userid uint64, followuid uint64) (bool, error) {
	_, err := db.c.Exec("INSERT INTO follow (uid, fuid) VALUES (?, ?)", userid, followuid)
	if err != nil {
		return false, err
	}
	return true, err
}
