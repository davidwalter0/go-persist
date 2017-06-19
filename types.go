package persist

import (
	"database/sql"
	"time"

	"html/template"
)

// Database driver info configure object
type Database struct {
	Driver   string        // database driver for backend, e.g. mysql, psql
	Host     string        // host of database
	Port     int           // port on which db
	Database string        // db name
	User     string        // db user
	Password string        // user password
	DropAll  bool          // re-initialize schema [for testing]
	Timeout  time.Duration //
	*sql.DB  `ignore:"true"`
}

type User struct {
	Id   int
	Name string
}

type Session struct {
	Id              string
	Authenticated   bool
	Unauthenticated bool
	User            User
}

type Comment struct {
	Id          int
	Name        string
	Email       string
	CommentText string
}

type Page struct {
	Id         int
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Comments   []Comment
	Session    Session
	GUID       string
}
