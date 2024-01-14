package database

import "errors"

// GetPostComments allows to get all the comments under a post.
// Request will fail if postid doesn't exist
func (db *appdbimpl) GetPostComments(postid uint64) ([]Comment, error) {
	const (
		commentQuery = "SELECT * FROM comment WHERE postid = ?"
	)

	// First check if post exist
	check, err := db.CheckPostByPostid(postid)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errors.New("post doesn't exist")
	}

	var comments []Comment
	rows, err := db.c.Query(commentQuery, postid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.Commentid, &comment.Message, &comment.Datetime, &comment.Postid, &comment.Userid)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, err

}
