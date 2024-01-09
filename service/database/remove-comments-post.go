package database

// RemoveCommentsFromPost allows to remove all comments under a post.
// Function will return nil for success, otherwise an error.
func (db *appdbimpl) RemoveCommentsFromPost(postid uint64) error {
	_, err := db.c.Exec("DELETE FROM comment WHERE postid = ?", postid)
	return err
}
