package database

// CheckCommentByCommentid checks if a specified comment exists.
// Request will return true if exists, otherwise false.
func (db *appdbimpl) CheckCommentByCommentid(commentid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM comment WHERE commentid = ?", commentid).Scan(&count)
	return count > 0, err
}
