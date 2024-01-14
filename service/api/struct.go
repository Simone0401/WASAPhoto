package api

import (
	"github.com/Simone0401/WASAPhoto/service/database"
	"regexp"
)

const (
	UserUsernameRegex   string = "^[A-Za-z0-9]{3,20}$"
	MessageCommentRegex string = "^[a-zA-Z0-9.,!?;:'\"\\s]+$"
)

// User struct represents a user in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package.
type User struct {
	Userid   uint64 `json:"user_id"`
	Username string `json:"username" validate:"min=3, max=20"`
}

// Comment struct represents a comment in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package.
type Comment struct {
	Commentid uint64 `json:"id"`
	Userid    uint64 `json:"uid"`
	Postid    uint64 `json:"postid"`
	Message   string `json:"message" validate:"min=1, max=256"`
	Datetime  string `json:"comment_datetime"`
}

// Post struct represents a post structure in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package.
type Post struct {
	Postid   uint64    `json:"postid"`
	Uid      uint64    `json:"uid"`
	Likes    uint64    `json:"likes" validate:"min=0"`
	Comments []Comment `json:"comments" validate:"dive"` // Validate Comments slice element, too
	Datetime string    `json:"upload_datetime" validate:"datetimeformat"`
}

// ProfileInfo struct represents a Profile structure in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package.
type ProfileInfo struct {
	User      User   `json:"user" validate:"dive"`
	Numpost   uint64 `json:"numpost" validate:"min=0"`
	Followers uint64 `json:"follower" validate:"min=0"`
	Following uint64 `json:"following" validate:"min=0"`
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

// FromDatabase populates the struct with data from the database, overwriting all values.
func (c *Comment) FromDatabase(comment database.Comment) error {
	c.Commentid = comment.Commentid
	c.Userid = comment.Userid
	c.Postid = comment.Postid
	c.Message = comment.Message
	c.Datetime = comment.Datetime
	return nil
}

// ToDatabase returns comment in a database-compatible representation
func (c *Comment) ToDatabase() database.Comment {
	return database.Comment{
		Commentid: c.Commentid,
		Userid:    c.Userid,
		Postid:    c.Postid,
		Message:   c.Message,
		Datetime:  c.Datetime,
	}
}

// FromDatabase populates the struct with data from the database, overwriting all values.
func (p *Post) FromDatabase(post database.Post) error {
	p.Postid = post.Postid
	p.Uid = post.Uid
	p.Likes = post.Likes
	for _, comment := range post.Comments {
		p.Comments = append(p.Comments, Comment(comment))
	}
	p.Datetime = post.Datetime
	return nil
}

// ToDatabase returns post in a database-compatible representation
func (p *Post) ToDatabase() database.Post {
	var postDatabase database.Post
	postDatabase.Postid = p.Postid
	postDatabase.Uid = p.Uid
	postDatabase.Likes = p.Likes
	for _, comment := range p.Comments {
		postDatabase.Comments = append(postDatabase.Comments, database.Comment(comment))
	}
	postDatabase.Datetime = p.Datetime
	return postDatabase
}

// FromDatabase populates the struct with data from the database, overwriting all values.
func (p *ProfileInfo) FromDatabase(profile database.Profile) error {
	err := p.User.FromDatabase(profile.User)
	if err != nil {
		return err
	}
	p.Numpost = profile.NumPost
	p.Followers = profile.Followers
	p.Following = profile.Following
	return nil
}

// ToDatabase returns profile in a database-compatible representation
func (p *ProfileInfo) ToDatabase() database.Profile {
	var profileDatabase database.Profile
	profileDatabase.User = p.User.ToDatabase()
	profileDatabase.NumPost = p.Numpost
	profileDatabase.Followers = p.Followers
	profileDatabase.Following = p.Following
	return profileDatabase
}

// IsValid checks the validity of the content. In particular, username should be in its range of validity.
// Note that the ID is not checked.
func (u *User) IsValid() bool {
	regexPattern := regexp.MustCompile(UserUsernameRegex)
	return regexPattern.MatchString(u.Username) && len(u.Username) > 2 && len(u.Username) < 21
}

// IsValid checks the validity of the content. In particular, comment message should be in it range of validity.
// Note that ID is not checked.
func (c *Comment) IsValid() bool {
	regexPattern := regexp.MustCompile(MessageCommentRegex)
	return regexPattern.MatchString(c.Message) && len(c.Message) > 0 && len(c.Message) < 257
}
