package database

// RemoveLikesFromPost allows to remove all likes under a post.
// Function will return nil for success, otherwise an error.
func (db *appdbimpl) RemoveLikesFromPost(postid uint64) error {
	_, err := db.c.Exec("DELETE FROM like WHERE postid = ?", postid)
	return err
}
