package mysql

import (
	uuid "github.com/davidwalter0/go-persist/uuid"

	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var DropAll bool

func Reinitialize() {
	DropAll = true
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func Connect(driver, db string) *sql.DB {
	DB, err := sql.Open(driver, db)
	checkErr(err)
	return DB
}

// Schema given that permissions grant table configuration, initialize the
// current database tables for this project.
var Schema map[string][]string = map[string][]string{
	"pages": []string{
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
	"comments": []string{
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
	"users": []string{

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
	"sessions": []string{
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

func Initialize(db *sql.DB, Schema map[string][]string) {
	// fmt.Println(db.QueryRow(`DROP TABLE pages cascade;`))
	for table, schema := range Schema {
		// fmt.Println(db.QueryRow(`DROP TABLE pages cascade;`))
		if DropAll {
			fmt.Println(db.QueryRow(fmt.Sprintf(`DROP TABLE %s cascade;`, table)))
		}
		for _, scheme := range schema {
			fmt.Printf("\ndb.QueryRow: %s:\n%s\n", table, scheme)
			fmt.Println(db.QueryRow(scheme))
		}
		// fmt.Println(db.Exec(table))
	}
	// for _, table := range Schema {
	// 	fmt.Printf("\ndb.QueryRow:\n%s\n", table)
	// 	fmt.Println(db.QueryRow(table))
	// }
	// 	fmt.Println(db.QueryRow(`
	// CREATE TABLE pages (
	//   id int(11) unsigned NOT NULL AUTO_INCREMENT,
	//   page_guid varchar(256) NOT NULL DEFAULT '',
	//   page_title varchar(256) DEFAULT NULL,
	//   page_content mediumtext,
	//   page_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	//   PRIMARY KEY (id),
	//   UNIQUE KEY page_guid (page_guid)
	// ) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
	// `))
	guid := uuid.GUID().String()
	fmt.Println("guid", guid)
	row := db.QueryRow("INSERT INTO pages (page_guid, page_title, page_content, page_date) VALUES ('" + guid + "', 'Hello, World', 'I''m so glad you found this page!  It''s been sitting patiently on the Internet for some time, just waiting for a visitor.', CURRENT_TIMESTAMP)")
	fmt.Printf("%v\n", *row)
	// fmt.Println("another one")

	rows, err := db.Query("select * from pages")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var page_id int
		var page_guid string
		var page_title string
		var page_content string
		var page_date time.Time
		rows.Scan(&page_id, &page_guid, &page_title, &page_content, &page_date)
		fmt.Println(page_id, page_guid, page_title, page_content, page_date)
	}

}
