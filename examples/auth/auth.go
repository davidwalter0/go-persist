package main

import (
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/schema"
	"github.com/davidwalter0/go-persist/uuid"
	"github.com/davidwalter0/twofactor"
)

var authDB = &persist.Database{}
var standAlone bool

func init() {
	if standAlone {
		authDB.ConfigEnvWPrefix("AUTH", false)
		authDB.Connect()
		authDB.DropAll(AuthSchema)
		authDB.Initialize(AuthSchema)
		log.Println(*authDB)
	}
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

// AuthKey accessible object for database authentication table I/O
type AuthKey struct {
	Email  string `json:"email"`
	Issuer string `json:"issuer"`
}

// Auth object db I/O for auth table
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
	db      *persist.Database `ignore:"true"`
	totp    []byte            `ignore:"true"` // raw otp bytes
	key     []byte            `ignore:"true"` // raw key bytes
	png     []byte            `ignore:"true"` // PNG byte array
	otp     *twofactor.Totp   `ignore:"true"` // OTP object
}

// NewKey create the key fields for an auth struct, notice that email
// uses account
func NewKey(email, issuer string) *Auth {
	return &Auth{
		Email:  email,
		Issuer: issuer,
		db:     authDB,
	}
}

// NewAuth initialize an auth struct, notice that email uses account
func NewAuth(email, issuer, hash string, key, totpBytes []byte, digits int) *Auth {
	return &Auth{
		Email:  email,
		Issuer: issuer,
		Hash:   "sha1",
		Digits: digits,
		Totp:   base64.StdEncoding.EncodeToString(totpBytes),
		Key:    base32.StdEncoding.EncodeToString(key),
		totp:   totpBytes,
		db:     authDB,
	}
}

// CopyAuth initialize an auth struct, notice that email uses account
func (auth *Auth) CopyAuth(from *Auth) {
	auth.Email = from.Email
	auth.Issuer = from.Issuer
	auth.Hash = from.Hash
	auth.Digits = from.Digits
	auth.Totp = from.Totp
	auth.Key = from.Key
	auth.totp = from.totp
	auth.key = from.key
	auth.otp = from.otp
	auth.png = from.png
	auth.db = from.db
}

// CopyKey initialize an auth struct, notice that email uses account
func (auth *Auth) CopyKey(from *Auth) {
	auth.Email = from.Email
	auth.Issuer = from.Issuer
	auth.db = from.db
}

// Create a row in a table
func (auth *Auth) Create() {
	authDB := auth.db
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

// Read row from db using auth key fields for query
func (auth *Auth) Read() bool {
	log.Println(*auth)
	authDB := auth.db
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
	count := auth.Count()
	fmt.Println("Count", count)
	return count != 0
}

// Update row from db using auth key fields
func (auth *Auth) Update() {
	authDB := auth.db
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

// Delete row from db using auth key fields
func (auth *Auth) Delete() {
	authDB := auth.db
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
	authDB := auth.db
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
