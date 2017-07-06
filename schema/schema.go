package schema

type SchemaText []string
type SchemaConfig struct {
	Table  string     `json:"table"`
	Schema SchemaText `json:"schema"`
}

type DBSchema map[string]SchemaText
