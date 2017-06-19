package uuid

import (
	// "fmt"
	"github.com/google/uuid"

	// "github.com/twinj/uuid"
	"regexp"
)

// A UUID representation compliant with specification in
// RFC 4122 document.
// type UUID [16]byte
const (
	UUID_SIZE = 16
)

// type UUID [UUID_SIZE]byte

// func (u *uuid.UUID) String() string {
// 	return string(u[:])
// }

func GUID() uuid.UUID {
	// var x = uuid.NewV1()
	// var b = x.Bytes()
	// var u UUID
	// for i := 0; i < UUID_SIZE && i < len(b); i++ {
	// 	u[i] = b[i]
	// }
	// fmt.Println(u.String())
	// fmt.Println(b)
	// fmt.Println(x)
	// fmt.Println(x.String())
	u, _ := uuid.NewUUID()
	return u
}

// Pattern used to parse hex string representation of the UUID.
// FIXME: do something to consider both brackets at one time,
// current one allows to parse string with only one opening
// or closing bracket.
const hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-" +
	"([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"

var re = regexp.MustCompile(hexPattern)

// a UUID object from a hex string representation.
// Function accepts UUID string in following formats:
//
//     "6ba7b814-9dad-11d1-80b4-00c04fd430c8"
//     "{6ba7b814-9dad-11d1-80b4-00c04fd430c8}"
//     "urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8"
//

func Validate(u uuid.UUID) (rc bool) {
	md := re.FindStringSubmatch(u.String())
	if md != nil {
		rc = true
	}
	return
}

// func Validate(u UUID) (rc bool) {
// 	md := re.FindStringSubmatch(u.String())
// 	if md != nil {
// 		rc = true
// 	}
// 	return
// }
