/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUsername(uid uint64) (string, error)
	SetUsername(uid uint64, name string) error
	GetUserByID(uid uint64) (User, error)
	GetUserByUsername(username string) (User, error)
	CheckExistsByUsername(username string) (bool, error)
	CreateUser(username string) (User, error)
	HasFollowed(userid uint64, followuid uint64) (bool, error)
	HasBanned(userid uint64, banneduid uint64) (bool, error)
	FollowUser(userid uint64, banneduid uint64) (bool, error)
	UnfollowUser(userid uint64, followuid uint64) (bool, error)
	HasMuted(userid uint64, muteduid uint64) (bool, error)
	BanUser(userid uint64, muteduid uint64) (bool, error)
	UnbanUser(userid uint64, muteduid uint64) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// User struct represent a user in every API call between this package and the outside world.
// Note that the internal representation of user in the database might be different.
type User struct {
	Userid   uint64
	Username string `validate:"min=3, max=20"`
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// check if table User exists
	err := checkTableUser(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure for table user: %w", err)
	}
	// check if table Post exists
	err = checkTablePost(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure for table post: %w", err)
	}
	// check if table Comment exists
	err = checkTableComment(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure for table comment: %w", err)
	}
	// check if table Like exists
	err = checkTableLike(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure for table like: %w", err)
	}
	// check if table Follow exists
	err = checkTableFollow(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure for table follow: %w", err)
	}
	// check if table Ban exists
	err = checkTableBan(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure for table ban: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

/*
 * checkTableUser check if User table already exists. If not exists, it will create that.
 */
func checkTableUser(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "CREATE TABLE user " +
			"(uid INTEGER PRIMARY KEY, " +
			"username TEXT NOT NULL CHECK(length(username) <= 20) UNIQUE);"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}

/*
 * checkTablePost check if Post table already exists. If not exists, it will create that.
 */
func checkTablePost(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='post';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "CREATE TABLE post " +
			"(postid INTEGER PRIMARY KEY, " +
			"message TEXT CHECK(length(message) <= 265), " +
			"uid INTEGER NOT NULL, " +
			"FOREIGN KEY (uid) REFERENCES user(uid))"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}

/*
 * checkTableComment check if Comment table already exists. If not exists, it will create that.
 */
func checkTableComment(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comment';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "CREATE TABLE comment " +
			"(commentid INTEGER PRIMARY KEY, " +
			"imageid INTEGER NOT NULL UNIQUE, " +
			"timestamp DATETIME, " +
			"postid INTEGER NOT NULL, " +
			"uid INTEGER NOT NULL, " +
			"FOREIGN KEY (uid) REFERENCES user(uid), " +
			"FOREIGN KEY (postid) REFERENCES post(postid))"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}

/*
 * checkTableLike check if Like table already exists. If not exists, it will create that.
 */
func checkTableLike(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='like';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "CREATE TABLE like " +
			"(uid INTEGER NOT NULL, " +
			"postid INTEGER NOT NULL, " +
			"PRIMARY KEY (uid, postid), " +
			"FOREIGN KEY (uid) REFERENCES user(uid), " +
			"FOREIGN KEY (postid) REFERENCES post(postid))"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}

/*
 * checkTableFollow check if Follow table already exists. If not exists, it will create that.
 */
func checkTableFollow(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='follow';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "CREATE TABLE follow " +
			"(uid INTEGER NOT NULL, " +
			"fuid INTEGER NOT NULL, " +
			"PRIMARY KEY (uid, fuid), " +
			"FOREIGN KEY (uid) REFERENCES user(uid), " +
			"FOREIGN KEY (fuid) REFERENCES user(uid))"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}

/*
 * checkTableBan check if Follow table already exists. If not exists, it will create that.
 */
func checkTableBan(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='ban';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "CREATE TABLE ban " +
			"(uid INTEGER NOT NULL, " +
			"buid INTEGER NOT NULL, " +
			"PRIMARY KEY (uid, buid), " +
			"FOREIGN KEY (uid) REFERENCES user(uid), " +
			"FOREIGN KEY (buid) REFERENCES user(uid))"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}
