package mock

import (
	"github.com/google/uuid"
	"time"
)

const (
	TestUUID      = "94310d8d-72d6-492a-bc30-27584627edb1"
	TestTimestamp = "2022-01-26T12:27:04.783936865Z"
)

// Generator is responsible for randomly generating new strings and tokens
// that might need to be mocked out to produce consistent output for tests
type Generator struct{}

// PSK returns a new random array of 16 bytes
func (g *Generator) PSK() ([]byte, error) {
	return nil, nil
}

// UUID generates a new V4 UUID
func (g *Generator) UUID() (uuid.UUID, error) {
	return uuid.Parse(TestUUID)
}

// Timestamp generates a timestamp of the current time
func (g *Generator) Timestamp() time.Time {
	t, _ := time.Parse(time.RFC3339, TestTimestamp)
	return t
}
