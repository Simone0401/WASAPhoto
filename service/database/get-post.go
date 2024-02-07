package database

import (
	"errors"
)

// GetPost allows to get all the information related to a post.
// Request will fail if postid doesn't exist
func (db *appdbimpl) GetPost(postid uint64) (Post, error) {
	const (
		postQuery = "SELECT * FROM post WHERE postid = ?"
	)

	// First check if post exist
	check, err := db.CheckPostByPostid(postid)
	if err != nil {
		return Post{}, err
	}
	if !check {
		return Post{}, errors.New("post doesn't exist")
	}

	var postDB Post
	err = db.c.QueryRow(postQuery, postid).Scan(&postDB.Postid, &postDB.Uid, &postDB.Datetime)
	if err != nil {
		return Post{}, err
	}

	likes, err := db.GetPostLikes(postid)
	if err != nil {
		return Post{}, err
	}
	postDB.Likes = uint64(len(likes))

	comments, err := db.GetPostComments(postid)
	if err != nil {
		return Post{}, err
	}
	postDB.Comments = comments

	return postDB, nil
}
