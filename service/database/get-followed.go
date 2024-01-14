package database

// GetFollowed allows to get a []uint64 followed ids for a specified user.
// Request will fail if uid doesn't exist
func (db *appdbimpl) GetFollowed(uid uint64) ([]uint64, error) {
	const (
		getFollowed = "SELECT follow.fuid from follow WHERE follow.uid = ?"
	)

	var fuids []uint64
	rows, err := db.c.Query(getFollowed, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read each row returned
	for rows.Next() {
		var fuid uint64
		if err = rows.Scan(&fuid); err != nil {
			return fuids, err
		}
		fuids = append(fuids, fuid)
	}

	return fuids, err
}
