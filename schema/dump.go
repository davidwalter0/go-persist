package schema

import (
	"fmt"
	"io"
)

func (schema *DBSchema) Dump(w io.Writer) {
	Dump(w, *schema)
}

func Dump(w io.Writer, Schema DBSchema) {
	for table, schema := range Schema {
		fmt.Fprintf(w, "- table: %s\n", table)
		fmt.Fprintf(w, "  schema:\n")
		for i, text := range schema {
			schema[i] = fmt.Sprintf(fmt.Sprintf("%s\n", text))
			fmt.Fprintf(w, "  - |-\n    %s", schema[i])
		}
	}
}

func (Schema *DBSchema) String() (text string) {
	for table, schema := range *Schema {
		text += fmt.Sprintf("- table: %s\n", table)
		text += fmt.Sprintf("  schema:")
		for _, definition := range schema {
			text += fmt.Sprintf("\n  - |-\n    %s", definition)
		}
	}
	return
}

func (schema *SchemaText) String(table string) string {
	return String(table, *schema)
}

func String(table string, schema SchemaText) (text string) {

	text += fmt.Sprintf("- table: %s\n", table)
	text += fmt.Sprintf("  schema:\n")
	for i, text := range schema {
		schema[i] = fmt.Sprintf("%s\n", text)
		text += fmt.Sprintf("\n  - |-    %s", schema[i])
	}

	return
}
