// go run dump-schema.go | tee yaml

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/schema"
	"github.com/davidwalter0/go-persist/uuid"
)

func Escape(text string) string {
	return strings.Replace(text, "'", "''", -1)
}

var Schema schema.DBSchema = schema.DBSchema{
	"schema": schema.SchemaText{
		`CREATE TABLE schema (
       schema_id  serial primary key,
       schema_guid varchar(256) NOT NULL unique,
       schema_database varchar(256) NOT NULL unique,
       schema_text text,
       schema_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       schema_changed timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`,
		`CREATE OR REPLACE FUNCTION update_schema_created_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.schema_changed = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql'`,
		`CREATE TRIGGER update_schema_ab_changetimestamp 
       BEFORE UPDATE ON schema 
       FOR EACH ROW EXECUTE PROCEDURE update_schema_created_column()`,
	},
}

func init() {
}

func main() {
	if false {
		Schema.Dump(os.Stdout)
	}
	var schemaDB = &persist.Database{}
	log.Println(*schemaDB)
	schemaDB.ConfigEnvWPrefix("SCHEMA", false)
	schemaDB.Connect()
	log.Println(*schemaDB)
	format := "Database: %v\nPort: %d\nUser: %s\nPassword: %s\nHost: %s\n"
	_, err := fmt.Printf(format, schemaDB.Database, schemaDB.Port, schemaDB.User, schemaDB.Password, schemaDB.Host)
	if err != nil {
		log.Fatal(err.Error())
	}
	schemaDB.DropAll(Schema)
	schemaDB.Initialize(Schema)

	var database = "schema"
	// 	for table, schema := range Schema {
	// 		fmt.Println("writing database", table)
	// 		INSERT := fmt.Sprintf(`
	// INSERT INTO schema (schema_guid, schema_database, schema_database, schema_text, schema_created)
	// VALUES ('%s', '%s', '%s', '%s', CURRENT_TIMESTAMP)`,
	// 			uuid.GUID().String(),
	// 			database+":"+uuid.GUID().String(), // table
	// 			// database,
	// 			// table+":"+uuid.GUID().String(), // table
	// 			schema.String(),
	// 		)
	// 		// fmt.Println(schemaDB.Exec(INSERT))
	// 		row := schemaDB.QueryRow(INSERT)
	// 		fmt.Printf("%v\n", *row)
	// 	}
	fmt.Println(Schema.String())
	fmt.Println("writing database", database)
	insert := fmt.Sprintf(`
INSERT INTO schema 
( schema_guid, 
  schema_database,
  schema_text,
  schema_created) 
VALUES ('%s', '%s', '%s', CURRENT_TIMESTAMP)`,
		uuid.GUID().String(),
		database+":"+uuid.GUID().String(), // table
		// database,
		// table+":"+uuid.GUID().String(), // table
		Escape(Schema.String()),
	)
	fmt.Println(schemaDB.Exec(insert))
	// row := schemaDB.QueryRow(insert)
	// fmt.Printf("%v\n", *row)

	// postgresql array in go
	// https://gist.github.com/adharris/4163702
	// rows := schemaDB.Query("SELECT schema_id, schema_guid, schema_database, schema_table, UNNEST(schema_text), schema_created, schema_changed FROM schema")

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

}

// // Parse the output string from the array type.
// // Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
// func parseArray(array string) []string {
// 	results := make([]string, 0)
// 	matches := arrayExp.FindAllStringSubmatch(array, -1)
// 	for _, match := range matches {
// 		s := match[valueIndex]
// 		// the string _might_ be wrapped in quotes, so trim them:
// 		s = strings.Trim(s, "\"")
// 		results = append(results, s)
// 	}
// 	return results
// }

// type StringSlice []string

// func (s *StringSlice) Scan(src interface{}) error {
// 	asBytes, ok := src.([]byte)
// 	if !ok {
// 		return error(errors.New("Scan source was not []bytes"))
// 	}

// 	asString := string(asBytes)
// 	fmt.Println("again", asString)
// 	parsed := parseArray(asString)
// 	(*s) = StringSlice(parsed)
// 	// (*s) = StringSlice(asString)

// 	return nil
// }

// var (
// 	// unquoted array values must not contain: (" , \ { } whitespace NULL)
// 	// and must be at least one char
// 	unquotedChar  = `[^",\\{}\s(NULL)]`
// 	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

// 	// quoted array values are surrounded by double quotes, can be any
// 	// character except " or \, which must be backslash escaped:
// 	quotedChar  = `[^"\\]|\\"|\\\\`
// 	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

// 	// an array value may be either quoted or unquoted:
// 	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

// 	// Array values are separated with a comma IF there is more than one value:
// 	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

// 	valueIndex int
// )
