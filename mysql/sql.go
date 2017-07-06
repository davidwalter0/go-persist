package mysql

import (
	"database/sql"
	schema "github.com/davidwalter0/go-persist/schema"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
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
       id int(11) unsigned NOT NULL AUTO_INCREMENT,
       page_guid varchar(256) NOT NULL DEFAULT '',
       page_title varchar(256) DEFAULT NULL,
       page_content mediumtext,
       page_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       PRIMARY KEY (id),
       UNIQUE KEY page_guid (page_guid)
       ) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1`,
	},
	"comments": schema.SchemaText{
		`CREATE TABLE comments (
       id int(11) unsigned NOT NULL AUTO_INCREMENT,
       page_id int(11) NOT NULL,
       comment_guid varchar(256) DEFAULT NULL,
       comment_name varchar(64) DEFAULT NULL,
       comment_email varchar(128) DEFAULT NULL,
       comment_text mediumtext,
       comment_date timestamp NULL DEFAULT NULL,
       PRIMARY KEY (id),
       KEY page_id (page_id)
       ) ENGINE=InnoDB DEFAULT CHARSET=latin1`,
	},
	"users": schema.SchemaText{

		`CREATE TABLE users (
       id int(11) unsigned NOT NULL AUTO_INCREMENT,
       user_name varchar(32) NOT NULL DEFAULT '',
       user_guid varchar(256) NOT NULL DEFAULT '',
       user_email varchar(128) NOT NULL DEFAULT '',
       user_password varchar(128) NOT NULL DEFAULT '',
       user_salt varchar(128) NOT NULL DEFAULT '',
       user_joined_timestamp timestamp NULL DEFAULT NULL,
       PRIMARY KEY (id)
     ) ENGINE=InnoDB DEFAULT CHARSET=latin1;`,
	},
	"sessions": schema.SchemaText{
		`CREATE TABLE sessions (
       id int(11) unsigned NOT NULL AUTO_INCREMENT,
       session_id varchar(256) NOT NULL DEFAULT '',
       user_id int(11) DEFAULT NULL,
       session_start timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       session_update timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
       session_active tinyint(1) NOT NULL,
       PRIMARY KEY (id),
       UNIQUE KEY session_id (session_id)
     ) ENGINE=InnoDB DEFAULT CHARSET=latin1`,
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
	// fmt.Println(db.QueryRow(`DROP TABLE pages cascade;`))
	for table, schema := range Schema {
		// fmt.Println(db.QueryRow(`DROP TABLE pages cascade;`))
		for _, scheme := range schema {
			fmt.Printf("\ndb.QueryRow: %s:\n%s\n", table, scheme)
			fmt.Println(db.QueryRow(scheme))
		}
	}
}
