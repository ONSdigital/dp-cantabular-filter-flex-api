package mongodb

import (
	"github.com/google/uuid"
)

type generator interface {
	UUID() (uuid.UUID, error)
}

type hasher interface {
	Hash([]byte) (string, error)
}
