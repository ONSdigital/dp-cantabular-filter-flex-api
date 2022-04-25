package mongodb

import (
	"time"

	"github.com/google/uuid"
)

type generator interface {
	UUID() (uuid.UUID, error)
	Timestamp() time.Time
}

type hasher interface {
	Hash([]byte) (string, error)
}
