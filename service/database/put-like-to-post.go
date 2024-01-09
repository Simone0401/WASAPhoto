package database

// LikePost allows to put a like from uid to a specified post.
// Function will return nil if no errors are present, an error otherwise.
func (db *appdbimpl) LikePost(postid uint64, userid uint64) error {
	var count uint64

	// Check if like is already put
	err := db.c.QueryRow("SELECT COUNT(*) FROM like WHERE uid = ? AND postid = ?", userid, postid).Scan(&count)
	if err != nil {
		return err
	}

	// Like already put
	if count > 0 {
		return nil
	}

	// Like have to put on
	_, err = db.c.Exec("INSERT INTO like(uid, postid) VALUES (?, ?)", userid, postid)

	return err
}
