package database

// GetFollowers allows to get a []uint64 followers ids for a specified user.
// Request will fail if uid doesn't exist
func (db *appdbimpl) GetFollowers(uid uint64) ([]uint64, error) {
	const (
		getFollowers = "SELECT follow.uid from follow WHERE follow.fuid = ?"
	)

	var fuids []uint64
	rows, err := db.c.Query(getFollowers, uid)
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
