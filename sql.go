package persist

import (
	"database/sql"

	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-persist/mysql"
	"github.com/davidwalter0/go-persist/pgsql"
	"github.com/davidwalter0/go-persist/schema"

	"fmt"

	"log"

	"encoding/json"
	"io/ioutil"
)

// debugging is not a secure property, and may write insecure
// information
var debugging bool

var drivers = []string{
	"mysql",
	"postgres",
	"pgsql",
}

// Connect initialize the driver and connect to the database
func Connect() *Database {
	db := &Database{}
	db.Connect()
	return db
}

// DriverName of the current driver
func (db *Database) DriverName() string {
	return db.Driver
}

// DriverNames array of available drivers
func (Database) DriverNames() []string {
	return drivers
}

// Handle defines the handle actions
type Handle interface {
	Connect() *Handle
	Close()
}

// CheckError standardize error handling
func CheckError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

// Configure the Database struct from the environment or flags
func (db *Database) Configure() {
	if err := cfg.Process("SQL", db); err != nil {
		log.Println(err)
	}

	if debugging { // insecure option
		var jsonText []byte
		jsonText, _ = json.MarshalIndent(db, "", "  ")
		_ = ioutil.WriteFile("Configure.SQL.json", jsonText, 0777)
	}
}

// ConfigEnvWPrefix fill an object with environment vars, if last call
// generate flags
func (db *Database) ConfigEnvWPrefix(envPrefix string, freeze bool) {
	if freeze {
		if err := cfg.Process(envPrefix, db); err != nil {
			log.Println(err)
		}
	} else {
		if err := cfg.ProcessHoldFlags(envPrefix, db); err != nil {
			log.Println(err)
		}
	}
	if debugging { // insecure option
		var jsonText []byte
		jsonText, _ = json.MarshalIndent(db, "", "  ")
		_ = ioutil.WriteFile("ConfigEnvWPrefix."+envPrefix+".json", jsonText, 0777)
	}
}

// ConnectString returns the db driver connection protocol string from
// the configured Database struct
func (db *Database) ConnectString() (text string) {

	switch db.DriverName() {
	case "pgsql", "postgres":
		text = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", db.User, db.Password, db.Database, db.Host, db.Port)
	case "mysql":
		text = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db.User, db.Password, db.Host, db.Port, db.Database)
	default:
		panic(fmt.Sprintf("ConnectString: driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
	}
	return
}

// Connect to the backend for the Database driver specified
func (db *Database) Connect() *Database {
	switch db.DriverName() {
	case "pgsql", "postgres":
		if debugging {
			fmt.Println(">", db.ConnectString())
		}
		db.DB = pgsql.Connect(db.DriverName(), db.ConnectString())
	case "mysql":
		if debugging {
			fmt.Println(">", db.ConnectString())
		}
		db.DB = mysql.Connect(db.DriverName(), db.ConnectString())
	default:
		panic(fmt.Sprintf("Connect: driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
	}
	if db.DB == nil {
		panic(fmt.Sprintf("Connect failed tried with driver %s valid drivers %v", db.DriverName(), db.DriverNames()))
	}
	return db
}

// Initialize a database from a schema definition iterating over the
// schema definition for each table and definition
func (db *Database) Initialize(schema schema.DBSchema) {
	switch db.DriverName() {
	case "pgsql", "postgres":
		pgsql.Initialize(db.DB, schema)
		return
	case "mysql":
		mysql.Initialize(db.DB, schema)
		return
	default:
		panic(fmt.Sprintf("Connect: driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
	}
	return
}

// Close the Database connection
func (db *Database) Close() error {
	return db.DB.Close()
}

// Insert a row to a database with optional arguments
func (db *Database) Insert(insert string, args ...interface{}) *sql.Row {
	row := db.DB.QueryRow(insert, args...)
	return row
}

// Query rows from a database
func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

// Prepare a query statement object
func (db *Database) Prepare(prepare string) (*sql.Stmt, error) {
	return db.DB.Prepare(prepare)
}

// DropAll remove the tables in this schema
func (db *Database) DropAll(Schema schema.DBSchema) {
	switch db.DriverName() {
	case "pgsql", "postgres":
		pgsql.DropAll(db.DB, Schema)
		return
	case "mysql":
		mysql.DropAll(db.DB, Schema)
		return
	default:
		panic(fmt.Sprintf("DropAll tables: required driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
	}

}
