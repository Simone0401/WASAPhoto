package database

import (
	"database/sql"
)

func (db *appdbimpl) SearchUserByUsername(username string) ([]User, error) {
	const (
		searchUseranameQuery = "SELECT * FROM user WHERE username LIKE ?"
	)

	rows, err := db.c.Query(searchUseranameQuery, username+"%")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Userid, &user.Username)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return users, rows.Err()
	}

	return users, nil
}
