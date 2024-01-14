package database

// GetProfilePosts allows to get profile Posts stream passing his uid.
// Request will fail if uid doesn't exist
func (db *appdbimpl) GetProfilePosts(uid uint64) ([]Post, error) {
	const (
		postsQuery = "SELECT post.postid, post.uid, post.timestamp FROM post WHERE post.uid = ? ORDER BY post.timestamp DESC"
	)

	var posts []Post

	// Make the query
	rows, err := db.c.Query(postsQuery, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Postid, &post.Uid, &post.Datetime)
		if err != nil {
			return nil, err
		}

		// Post correctly readed, get number of likes
		likes, err := db.GetPostLikes(post.Postid)
		if err != nil {
			return posts, err
		}
		post.Likes = uint64(len(likes))

		// Get comments
		post.Comments, err = db.GetPostComments(post.Postid)
		if err != nil {
			return posts, err
		}

		// Add post to the list
		posts = append(posts, post)
	}

	return posts, err

}
