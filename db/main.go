package main

import (
	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/pgsql"

	"fmt"
	"log"
)

func main() {
	var err error
	var db = &persist.Database{}

	db.ConfigEnvWPrefix("SQL", true)
	db.Connect()
	format := "Database: %v\nPort: %d\nUser: %s\nPassword: %s\nHost: %s\n"
	_, err = fmt.Printf(format, db.Database, db.Port, db.User, db.Password, db.Host)
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Initialize(pgsql.Schema)
}
