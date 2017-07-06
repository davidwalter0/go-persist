package persist

import (
	"database/sql"
	"time"

	"html/template"
)

// Database driver info configure object
type Database struct {
	Driver   string        `usage:"database driver for backend, e.g. mysql, psql"`
	Host     string        `usage:"host of database"`
	Port     int           `usage:"port on which database"`
	Database string        `usage:"database name"`
	User     string        `usage:"database user name"`
	Password string        `usage:"database user password"`
	Timeout  time.Duration `usage:"connect timeout"`
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
