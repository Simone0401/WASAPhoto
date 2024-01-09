package database

// UnlikePost allows to delete a user like to a specified post.
// Function will return nil if no errors are present, an error otherwise.
func (db *appdbimpl) UnlikePost(postid uint64, userid uint64) error {
	_, err := db.c.Exec("DELETE FROM like WHERE postid = ? AND uid = ?", postid, userid)
	return err
}
