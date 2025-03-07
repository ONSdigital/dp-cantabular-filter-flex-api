package mock

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TestUUID      = "94310d8d-72d6-492a-bc30-27584627edb1"
	TestTimestamp = "2022-01-26T12:27:04.783936865Z"
)

// Generator is responsible for generating mocked constant strings and tokens
// for tests
type Generator struct {
	URLHost string
}

// PSK returns a new constant array of 16 bytes
func (g *Generator) PSK() ([]byte, error) {
	return []byte("0123456789ABCDEF"), nil
}

// UUID generates a constant UUID
func (g *Generator) UUID() (uuid.UUID, error) {
	return uuid.Parse(TestUUID)
}

// Timestamp generates a constant timestamp
func (g *Generator) Timestamp() time.Time {
	t, _ := time.Parse(time.RFC3339, TestTimestamp)
	return t
}

// UniqueTimestamp generates a constant timestamp
func (g *Generator) UniqueTimestamp() primitive.Timestamp {
	t, _ := time.Parse(time.RFC3339, TestTimestamp)
	return primitive.Timestamp{
		//nolint:gosec // G115: integer overflow conversion int64 -> uint32 // acceptable until February 7, 2106
		T: uint32(t.Unix()),
		I: 1,
	}
}

// URL generates a URL from a constant host and a path made from a printf
// string + arguments
func (g *Generator) URL(_, path string, args ...interface{}) string {
	return g.URLHost + fmt.Sprintf(path, args...)
}
