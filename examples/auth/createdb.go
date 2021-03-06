package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/uuid"
)

func Escape(text string) string {
	return strings.Replace(text, "'", "''", -1)
}

func init() {
}

func main() {
	var format string
	var err error

	if false {
		AuthSchema.Dump(os.Stdout)
	}
	var authDB = &persist.Database{}
	var schemaDB = &persist.Database{}

	authDB.ConfigEnvWPrefix("AUTH", false)
	schemaDB.ConfigEnvWPrefix("SCHEMA", false)

	authDB.Connect()
	defer authDB.Close()

	authDB.DropAll(AuthSchema)
	authDB.Initialize(AuthSchema)

	format = "Database: %v\nPort: %d\nUser: %s\nPassword: %s\nHost: %s\n"
	_, err = fmt.Printf(format, authDB.Database, authDB.Port, authDB.User, authDB.Password, authDB.Host)
	if err != nil {
		log.Fatal(err.Error())
	}

	schemaDB.Connect()
	defer schemaDB.Close()

	format = "Database: %v\nPort: %d\nUser: %s\nPassword: %s\nHost: %s\n"
	_, err = fmt.Printf(format, schemaDB.Database, schemaDB.Port, schemaDB.User, schemaDB.Password, schemaDB.Host)
	if err != nil {
		log.Fatal(err.Error())
	}

	if true {
		var database = "auth"
		fmt.Println(AuthSchema.String())
		fmt.Println("writing database", database)
		insert := fmt.Sprintf(`
INSERT INTO schema 
( schema_guid, 
  schema_database,
  schema_text,
  schema_created
) 
VALUES ('%s', '%s', '%s', CURRENT_TIMESTAMP)`,
			uuid.GUID().String(),
			database,
			Escape(AuthSchema.String()),
		)
		fmt.Println(schemaDB.Exec(insert))
	}
	rows := schemaDB.Query("SELECT schema_id, schema_guid, schema_database, schema_text, schema_created, schema_changed FROM schema")
	defer rows.Close()

	for rows.Next() {
		var schemaID int
		var schemaGUID string
		var schemaDatabase string
		var schemaText string
		var schemaCreated time.Time
		var schemaChanged time.Time
		rows.Scan(
			&schemaID,
			&schemaGUID,
			&schemaDatabase,
			&schemaText,
			&schemaCreated,
			&schemaChanged)
		fmt.Println(
			schemaID,
			schemaGUID,
			schemaDatabase,
			schemaText,
			schemaCreated,
			schemaChanged)
	}

	auth := &Auth{Email: "walter.david@gmail.com", Issuer: "vpn0.me", Key: "fake key: 67b9cecb-6071-11e7-93b5-68f7284fe468", Totp: "fake totp: base64...", GUID: "67b9cecb-6071-11e7-93b5-68f7284fe468", Hash: "sha1", Digits: 6, db: authDB}
	check := &Auth{}
	deleter := &Auth{}
	check.CopyKey(auth)
	deleter.CopyKey(auth)

	auth.Create()
	fmt.Println("auth", auth)
	check.Read()
	fmt.Println(check)
	auth.Totp = "another fake totp"
	auth.Update()
	check.Read()
	fmt.Println("ok", check)

	deleter.Delete()
	deleter.Read()
	fmt.Println("ok", deleter)
}
