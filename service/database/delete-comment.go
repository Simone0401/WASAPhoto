package database

func (db *appdbimpl) DeleteComment(commentid uint64) error {
	_, err := db.c.Exec("DELETE FROM comment WHERE commentid = ?", commentid)
	return err
}
