package persist

import (
	"database/sql"
	"os"

	"github.com/davidwalter0/go-persist/mysql"
	"github.com/davidwalter0/go-persist/pgsql"

	"fmt"

	"github.com/davidwalter0/envflagstructconfig"
	"log"

	"encoding/json"
	"io/ioutil"
)

var drivers []string = []string{
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
		// panic(err)
	}
}

func (db *Database) Configure() {
	envflagstructconfig.Process("SQL", db)
	var jsonText []byte
	jsonText, _ = json.MarshalIndent(db, "", "  ")
	// fmt.Printf("\n%v\n", string(jsonText))
	ioutil.WriteFile("tmp.xyz.json", jsonText, 0777)
	// log.Printf("%v\n", db)
}

func (db *Database) ConnectString() (text string) {

	switch db.DriverName() {
	case "pgsql", "postgres":
		text = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", db.User, db.Password, db.Database, db.Host, db.Port)
	case "mysql":
		text = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db.User, db.Password, db.Host, db.Port, db.Database)
	default:
		panic(fmt.Sprintf("ConnectString: driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
		os.Exit(1)
	}
	return
}

func (db *Database) Connect() *Database {
	db.Configure()
	switch db.DriverName() {
	case "pgsql", "postgres":
		fmt.Println(">", db.ConnectString())
		db.DB = pgsql.Connect(db.DriverName(), db.ConnectString())
	case "mysql":
		fmt.Println(">", db.ConnectString())
		db.DB = mysql.Connect(db.DriverName(), db.ConnectString())
	default:
		panic(fmt.Sprintf("Connect: driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
		os.Exit(1)
	}
	return db
}

func (db *Database) Initialize() {
	switch db.DriverName() {
	case "pgsql", "postgres":
		pgsql.Reinitialize()
		pgsql.Initialize(db.DB, pgsql.Schema)
		return
	case "mysql":
		mysql.Reinitialize()
		mysql.Initialize(db.DB, mysql.Schema)
		return
	default:
		panic(fmt.Sprintf("Connect: driver mode unknown or empty %s valid drivers %v", db.DriverName(), db.DriverNames()))
	}
	return
}

func (db *Database) Close() {
	err := db.DB.Close()
	CheckError(err)
}

func (db *Database) Insert(insert string, args ...interface{}) *sql.Row {
	row := db.DB.QueryRow(insert, args...)
	return row
}

func (db *Database) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := db.DB.Query(query, args...)
	CheckError(err)
	return rows
}

func (db *Database) Prepare(prepare string) *sql.Stmt {
	statement, err := db.DB.Prepare(prepare)
	CheckError(err)
	return statement
}
