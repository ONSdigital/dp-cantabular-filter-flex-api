package mongodb

import (
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
)

// anyETag represents the wildchar that corresponds to not check the ETag value for update requests
const anyETag = "*"

// CreateEtag creates a new etag for when an update request is made.
func (c *Client) CreateEtag(current, update interface{}) (eTag string, err error) {
	b, err := bson.Marshal(update)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal bson: %w")
	}

	if h, ok := current.(hasher); ok {
		return h.Hash(b)
	}

	return "", nil
}
