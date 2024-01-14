package database

// GetProfileInfo allows to get user profile information passing an uid.
// Request will fail if uid doesn't exist
func (db *appdbimpl) GetProfileInfo(uid uint64) (Profile, error) {
	const (
		getNumberPostQuery = "SELECT COUNT(*) FROM post WHERE post.uid = ?"
	)

	var profile Profile
	var err error
	profile.User, err = db.GetUserByID(uid)
	if err != nil {
		return profile, err
	}

	err = db.c.QueryRow(getNumberPostQuery, uid).Scan(&profile.NumPost)
	if err != nil {
		return profile, err
	}

	followers, err := db.GetFollowers(uid)
	if err != nil {
		return profile, err
	}
	profile.Followers = uint64(len(followers))

	followed, err := db.GetFollowed(uid)
	if err != nil {
		return profile, err
	}
	profile.Following = uint64(len(followed))

	return profile, err
}
