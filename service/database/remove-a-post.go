package database

// RemovePost allows to remove a specified post if the specified user is the owner.
// Function will return nil if no errors are present, an error otherwise.
func (db *appdbimpl) RemovePost(postid uint64, userid uint64) error {
	_, err := db.c.Exec("DELETE FROM post WHERE postid = ? AND uid = ?", postid, userid)
	return err
}
