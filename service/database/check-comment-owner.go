package database

// CheckCommentOwner checks if a specified user is the owner of a specified comment.
// Request will return true if user is the owner, otherwise false.
func (db *appdbimpl) CheckCommentOwner(commentid uint64, userid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM comment WHERE commentid = ? AND uid = ?", commentid, userid).Scan(&count)
	return count > 0, err
}
