package generator

import (
	"crypto/rand"
	"github.com/google/uuid"
	"time"

	"github.com/pkg/errors"
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

// Timestamp generates a timestamp of the current time
func (g *Generator) Timestamp() time.Time {
	return time.Now()
}