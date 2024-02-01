package database

// CheckLike allows to check if a like is already put from uid to a specified post.
// Function will return nil if no errors are present, an error otherwise.
func (db *appdbimpl) CheckLike(postid uint64, userid uint64) (bool, error) {
	var count uint64

	// Check if like is already put
	err := db.c.QueryRow("SELECT COUNT(*) FROM like WHERE uid = ? AND postid = ?", userid, postid).Scan(&count)
	if err != nil {
		return false, err
	}

	// Like already put
	return count > 0, nil
}
