package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/schema"
	"github.com/davidwalter0/go-persist/uuid"
)

var authDB = &persist.Database{}

func init() {
	authDB.ConfigEnvWPrefix("AUTH", false)
	authDB.Connect()
	authDB.DropAll(AuthSchema)
	authDB.Initialize(AuthSchema)
	log.Println(*authDB)
}

// AuthSchema describes the table and triggers for persisting
// authentications from totp objects from twofactor
var AuthSchema = schema.DBSchema{
	"auth": schema.SchemaText{ // issuer <-> domain
		`CREATE TABLE auth (
       id  serial primary key,
       guid varchar(256) NOT NULL unique,
       email varchar(256) NOT NULL,
       issuer varchar(256) NOT NULL, 
       hash varchar(64) NOT NULL DEFAULT 'sha1', 
       digits integer NOT NULL DEFAULT 6, 
       key varchar(256) NOT NULL, 
       totp varchar(1024) NOT NULL, 
       created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       changed timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`,
		`CREATE UNIQUE INDEX auth_unique_idx on auth (email, issuer)`,
		`CREATE OR REPLACE FUNCTION update_created_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.changed = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql'`,
		`CREATE TRIGGER update_ab_changetimestamp 
       BEFORE UPDATE ON auth
       FOR EACH ROW EXECUTE PROCEDURE update_created_column()`,
	},
}

// Auth accessible object for database authentication table
// I/O
type AuthKey struct {
	Email  string `json:"email"`
	Issuer string `json:"issuer"`
}

type Auth struct {
	ID      int               `json:"id"`
	GUID    string            `json:"guid"`
	Email   string            `json:"email"`
	Issuer  string            `json:"issuer"`
	Hash    string            `json:"hash"`
	Digits  int               `json:"digits"`
	Created time.Time         `json:"created"`
	Changed time.Time         `json:"changed"`
	Key     string            `json:"key"  usage:"base32 encoded totp key"`
	Totp    string            `json:"totp" usage:"base64 encoded totp object"`
	DB      *persist.Database `ignore:"true"`
}

type DBIO interface {
	Create() error
	Read() error
	Update() error
	Delete() error
}

// NewAuth creates an auth object and initializes the
// connection object
func NewAuth(authDB *persist.Database) *Auth {
	return &Auth{DB: authDB}
}

func (auth *Auth) Create() {
	authDB := auth.DB
	// ignore DB & id
	insert := fmt.Sprintf(`
INSERT INTO auth 
( guid, 
  email,
  issuer,
  hash,
  digits,
  key,
  totp,
  created,
  changed
)
VALUES ('%s', '%s', '%s', '%s', %d, '%s', '%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		uuid.GUID().String(),
		auth.Email,
		auth.Issuer,
		auth.Hash,
		auth.Digits,
		auth.Key,
		auth.Totp,
	)
	fmt.Println(insert)
	fmt.Println(authDB.Exec(insert))
	fmt.Println("Count", auth.Count())
}

func (auth *Auth) Read() {
	authDB := auth.DB
	// ignore DB & id
	query := fmt.Sprintf(`
SELECT 
  id,
  guid, 
  email,
  issuer,
  hash,
  digits,
  key,
  totp,
  created,
  changed
FROM
   auth 
WHERE
  email = '%s'
AND
  issuer = '%s'
`,
		auth.Email,
		auth.Issuer,
	)
	fmt.Println(query)
	rows := authDB.Query(query)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&auth.ID,
			&auth.GUID,
			&auth.Email,
			&auth.Issuer,
			&auth.Hash,
			&auth.Digits,
			&auth.Key,
			&auth.Totp,
			&auth.Created,
			&auth.Changed); err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		fmt.Println(
			auth.ID,
			auth.GUID,
			auth.Email,
			auth.Issuer,
			auth.Hash,
			auth.Digits,
			auth.Key,
			auth.Totp,
			auth.Created,
			auth.Changed)
	}
	fmt.Println("Count", auth.Count())
}

func (auth *Auth) Update() {
	authDB := auth.DB
	// ignore DB & id
	update := fmt.Sprintf(`
UPDATE
  auth
SET
  hash    = '%s',
  digits  =  %d,
  key     = '%s',
  totp    = '%s'
WHERE
  email  = '%s'
AND
  issuer = '%s'
`,
		// set
		auth.Hash,
		auth.Digits,
		auth.Key,
		auth.Totp,
		// where
		auth.Email,
		auth.Issuer,
	)
	var err error
	var rows *sql.Rows
	var result sql.Result
	fmt.Println(update)
	result, err = authDB.Exec(update)
	fmt.Println("update result", result, "error", err)
	rows = authDB.Query("SELECT * FROM auth")
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&auth.ID,
			&auth.GUID,
			&auth.Email,
			&auth.Issuer,
			&auth.Hash,
			&auth.Digits,
			&auth.Key,
			&auth.Totp,
			&auth.Created,
			&auth.Changed); err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		fmt.Println(
			auth.ID,
			auth.GUID,
			auth.Email,
			auth.Issuer,
			auth.Hash,
			auth.Digits,
			auth.Key,
			auth.Totp,
			auth.Created,
			auth.Changed)
	}
	fmt.Println("Count", auth.Count())
}

func (auth *Auth) Delete() {
	authDB := auth.DB
	// ignore DB & id
	delete := fmt.Sprintf(`
DELETE FROM
  auth
WHERE
  email  = '%s'
AND
  issuer = '%s'
`,
		// where
		auth.Email,
		auth.Issuer,
	)
	var err error
	var rows *sql.Rows
	var result sql.Result
	fmt.Println(delete)
	result, err = authDB.Exec(delete)
	fmt.Println("delete result", result, "error", err)
	rows = authDB.Query("SELECT * FROM auth")
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&auth.ID,
			&auth.GUID,
			&auth.Email,
			&auth.Issuer,
			&auth.Hash,
			&auth.Digits,
			&auth.Key,
			&auth.Totp,
			&auth.Created,
			&auth.Changed); err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		fmt.Println(
			auth.ID,
			auth.GUID,
			auth.Email,
			auth.Issuer,
			auth.Hash,
			auth.Digits,
			auth.Key,
			auth.Totp,
			auth.Created,
			auth.Changed)
	}
	fmt.Println("Count", auth.Count())
}

// Count rows for keys in auth
func (auth *Auth) Count() (count int) {
	authDB := auth.DB
	query := fmt.Sprintf(`
SELECT
  COUNT(*) 
FROM
  auth
WHERE 
  email  = '%s'
AND
  issuer = '%s'
`,
		// where
		auth.Email,
		auth.Issuer,
	)

	row := authDB.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		log.Println("Row count query error", err)
	}
	return count
}
