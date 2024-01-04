package database

// AddPost allows to create a new post for a specific user.
// Function will return the created new post id .
func (db *appdbimpl) AddPost(userid uint64) (uint64, error) {
	var maxId uint64

	// RECOVER MAX(ID)
	err := db.c.QueryRow("SELECT MAX(postid) FROM post").Scan(&maxId)
	if err != nil {
		maxId = 0
	}

	postid := maxId + 1

	_, err = db.c.Exec("INSERT INTO post VALUES (?, ?, (SELECT datetime('now', '+1 hours')))", postid, userid)

	if err != nil {
		return 0, err
	}

	return postid, err
}
