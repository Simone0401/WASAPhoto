package database

// LikePost allows to put a like from uid to a specified post.
// Function will return nil if no errors are present, an error otherwise.
func (db *appdbimpl) LikePost(postid uint64, userid uint64) error {
	// Check if like is already put
	check, err := db.CheckLike(postid, userid)
	if err != nil {
		return err
	}

	// Like already put
	if check {
		return nil
	}

	// Like have to put on
	_, err = db.c.Exec("INSERT INTO like(uid, postid) VALUES (?, ?)", userid, postid)

	return err
}
