package database

// CheckPostByPostid checks if a post id is already present in database.
// Request will fail if post id doesn't exist
func (db *appdbimpl) CheckPostByPostid(postid uint64) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM post WHERE postid = ?", postid).Scan(&count)
	return count > 0, err
}
