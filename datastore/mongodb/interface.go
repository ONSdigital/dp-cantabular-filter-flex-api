package mongodb

import (
	"time"

	"github.com/google/uuid"
)

type generator interface {
	PSK() ([]byte, error)
	UUID() (uuid.UUID, error)
	Timestamp() time.Time
}

type hasher interface {
	Hash([]byte) (string, error)
}

type errNotFound interface {
	NotFound() bool
}

type errConflict interface {
	Conflict() bool
}

type errUnavailable interface {
	Unavailable() bool
}
