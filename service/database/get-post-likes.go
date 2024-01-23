package database

import (
	"database/sql"
)

// GetPostLikes allows to get a []uint64 follower ids that put like to a specified post.
// Request will fail if postid doesn't exist
func (db *appdbimpl) GetPostLikes(postid uint64) ([]uint64, error) {
	const (
		likesQuery = "SELECT like.uid FROM like WHERE like.postid = ?"
	)

	var uidlikes []uint64

	rows, err := db.c.Query(likesQuery, postid)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var uidlike uint64
		err = rows.Scan(&uidlike)
		if err != nil {
			return uidlikes, err
		}
		uidlikes = append(uidlikes, uidlike)
	}

	if rows.Err() != nil {
		return uidlikes, rows.Err()
	}

	return uidlikes, err

}
