package main

import (
	db "github.com/davidwalter0/go-persist"

	"fmt"
	"log"
)

func main() {
	var err error
	var db = &db.Database{}

	db.Connect()
	format := "Database: %v\nPort: %d\nUser: %s\nPassword: %s\nHost: %s\n"
	_, err = fmt.Printf(format, db.Database, db.Port, db.User, db.Password, db.Host)
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Initialize()
}
