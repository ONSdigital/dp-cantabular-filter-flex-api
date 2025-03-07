package generator

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Generator is responsible for randomly generating new strings and tokens
// that might need to be mocked out to produce consistent output for tests
type Generator struct{}

// New returns a new Generator
func New() *Generator {
	return &Generator{}
}

// PSK returns a new random array of 16 bytes
func (g *Generator) PSK() ([]byte, error) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return nil, errors.WithStack(err)
	}

	return key, nil
}

// UUID generates a new V4 UUID
func (g *Generator) UUID() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	return id, errors.WithStack(err)
}

// UniqueTimestamp generates a timestamp of the current time in the
// special format required by mongoDB
func (g *Generator) UniqueTimestamp() primitive.Timestamp {
	return primitive.Timestamp{
		//nolint:gosec // G115: integer overflow conversion int64 -> uint32 // acceptable until February 7, 2106
		T: uint32(time.Now().Unix()),
		I: 1,
	}
}

// Timestamp generates a timestamp of the current time
func (g *Generator) Timestamp() time.Time {
	return time.Now()
}

// URL generates a URL from a host and a path made from a printf string
// + arguments
func (g *Generator) URL(host, path string, args ...interface{}) string {
	return host + fmt.Sprintf(path, args...)
}
