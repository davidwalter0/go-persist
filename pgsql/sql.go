package pgsql

import (
	schema "github.com/davidwalter0/go-persist/schema"

	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

// Connect using driver and database name
func Connect(driver, db string) *sql.DB {
	DB, err := sql.Open(driver, db)
	checkErr(err)
	return DB
}

// Schema given that permissions grant table configuration, initialize the
// current database tables for this project.
var Schema = schema.DBSchema{
	"pages": schema.SchemaText{
		`CREATE TABLE pages (
       id  serial primary key,
       page_guid varchar(256) NOT NULL DEFAULT '' unique,
       page_title varchar(256) DEFAULT NULL,
       page_content text,
       page_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE OR REPLACE FUNCTION update_page_date_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.page_date = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql'`,
		`CREATE TRIGGER update_ab_changetimestamp 
       BEFORE UPDATE ON pages 
       FOR EACH ROW EXECUTE PROCEDURE update_page_date_column()`,
	},
	"comments": schema.SchemaText{
		`CREATE TABLE comments (
       id serial primary key,
       page_id int,
       comment_guid varchar(256) DEFAULT NULL,
       comment_name varchar(64) DEFAULT NULL,
       comment_email varchar(128) DEFAULT NULL,
       comment_text text,
       comment_date timestamp NULL DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE OR REPLACE FUNCTION update_comments_date_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.comment_date = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql'`,
		`CREATE TRIGGER update_comments_changetimestamp 
       BEFORE UPDATE ON comments
       FOR EACH ROW EXECUTE PROCEDURE update_page_date_column()`,
	},
	"users": schema.SchemaText{
		`CREATE TABLE users (
      id serial primary key,
      user_name varchar(32) NOT NULL DEFAULT '',
      user_guid varchar(256) NOT NULL DEFAULT '',
      user_email varchar(128) NOT NULL DEFAULT '',
      user_password varchar(128) NOT NULL DEFAULT '',
      user_salt varchar(128) NOT NULL DEFAULT '',
      user_joined_timestamp timestamp NULL DEFAULT NULL)`,
		`CREATE OR REPLACE FUNCTION update_comments_date_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.comment_date = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql';
  `,
		`CREATE TRIGGER update_comments_changetimestamp 
       BEFORE UPDATE ON comments
       FOR EACH ROW EXECUTE PROCEDURE update_comments_date_column()`,
	},
	"sessions": schema.SchemaText{
		`CREATE TABLE sessions (
       id serial primary key,
       session_id varchar(256) NOT NULL unique,
       user_id int DEFAULT NULL,
       session_start timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       session_update timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
       session_active int NOT NULL)`,
		`CREATE OR REPLACE FUNCTION update_sessions_start_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.session_start = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql'`,
		`CREATE TRIGGER update_sessions_changetimestamp 
       BEFORE UPDATE ON sessions
       FOR EACH ROW EXECUTE PROCEDURE update_sessions_start_column()`,
	},
}

// DropAll remove the tables in this schema
func DropAll(db *sql.DB, Schema schema.DBSchema) {
	for table, _ := range Schema {
		fmt.Println(db.QueryRow(fmt.Sprintf(`DROP TABLE %s cascade;`, table)))
	}
}

// Initialize a database from the given schema, the create operation
// is idempotent, and can be called multiple times without issues, if
// DropAll is false. Assumes that the uid running the process has
// given that permissions granted for table configuration, initialize
// the current database tables for this project.
func Initialize(db *sql.DB, Schema schema.DBSchema) {
	for table, schema := range Schema {
		// fmt.Println(db.QueryRow(`DROP TABLE pages cascade;`))
		for _, scheme := range schema {
			fmt.Printf("\ndb.QueryRow: %s:\n%s\n", table, scheme)
			fmt.Println(db.QueryRow(scheme))
		}
	}
}
