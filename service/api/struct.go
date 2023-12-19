package api

import (
	"github.com/Simone0401/WASAPhoto/service/database"
	"regexp"
)

const (
	UserUsernameRegex string = "^[A-Za-z0-9]{3,20}$"
)

// User struct represent a user in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package.
type User struct {
	Userid   uint64 `json:"user_id"`
	Username string `json:"username" validate:"min=3, max=20"`
}

// FromDatabase populates the struct with data from the database, overwriting all values.
// You might think this is code duplication, which is correct. However, it's "good" code duplication because it allows
// us to uncouple the database and API packages.
// Suppose we were using the "database.User" struct inside the API package; in that case, we were forced to conform
// either the API specifications to the database package or the other way around. However, very often, the database
// structure is different from the structure of the REST API.
// Also, in this way the database package is freely usable by other packages without the assumption that structs from
// the database should somehow be JSON-serializable (or, in general, serializable).
func (u *User) FromDatabase(user database.User) error {
	u.Userid = user.Userid
	u.Username = user.Username
	return nil
}

// ToDatabase returns the user in a database-compatible representation
func (u *User) ToDatabase() database.User {
	return database.User{
		Userid:   u.Userid,
		Username: u.Username,
	}
}

// IsValid checks the validity of the content. In particular, username should be in its range of validity.
// Note that the ID is not checked.
func (u *User) IsValid() bool {
	regexPattern := regexp.MustCompile(UserUsernameRegex)
	return regexPattern.MatchString(u.Username) && len(u.Username) > 2 && len(u.Username) < 21
}
