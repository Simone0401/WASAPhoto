package database

func (db *appdbimpl) AddComment(userid uint64, postid uint64, message string) (Comment, error) {
	var commentId uint64
	err := db.c.QueryRow("SELECT COUNT(commentid) FROM comment").Scan(&commentId)
	if err != nil {
		return Comment{}, err
	} else if commentId > 0 {
		err = db.c.QueryRow("SELECT MAX(commentid) FROM comment").Scan(&commentId)
		if err != nil {
			return Comment{}, err
		}
	} else {
		commentId = 0
	}

	commentId = commentId + 1
	_, err = db.c.Exec("INSERT INTO comment(commentid, message, timestamp, postid, uid) VALUES (?, ?, datetime('now', '+1 hours'), ?, ?)", commentId, message, postid, userid)
	if err != nil {
		return Comment{}, err
	}

	var comment Comment
	err = db.c.QueryRow("SELECT * FROM comment WHERE commentid = ?", commentId).Scan(&comment.Commentid, &comment.Message, &comment.Datetime, &comment.Postid, &comment.Userid)
	if err != nil {
		return Comment{}, err
	}
	return comment, err
}
