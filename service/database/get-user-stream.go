package database

import "strings"

// GetUserStream allows to get a user Posts stream passing his uid.
// Request will fail if uid doesn't exist
func (db *appdbimpl) GetUserStream(uid uint64) ([]Post, error) {
	const (
		postsQueryBase     = "SELECT post.postid, post.uid, post.timestamp FROM post WHERE post.uid IN "
		postOrderQueryBase = ") ORDER BY post.timestamp DESC"
	)

	// Get the list of followed
	followed, err := db.GetFollowed(uid)
	if err != nil {
		return nil, err
	}

	// Make placeholder string for IN query
	placeholders := make([]string, len(followed))
	values := make([]interface{}, len(followed))
	for i, follower := range followed {
		placeholders[i] = "?"
		values[i] = follower
	}

	// Add comma between every pair of '?'
	placeholdersStr := strings.Join(placeholders, ", ")

	// Build final query
	query := postsQueryBase + "(" + placeholdersStr + postOrderQueryBase

	rows, err := db.c.Query(query, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

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
