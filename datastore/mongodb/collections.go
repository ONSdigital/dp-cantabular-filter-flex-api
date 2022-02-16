package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-mongodb/v3/dplock"

	"github.com/pkg/errors"
)

// Collections holds information about the mongodb collections
// relevant to this service
type Collections struct {
	filters       *Collection
	filterOutputs *Collection
}

// Collection represents a collection in mongodb. Holds
// The name(s), lock client and any other relevant information.
type Collection struct {
	name       string
	lockClient *dplock.Lock
}

func (c *Collection) lock(ctx context.Context, id string) (string, error) {
	l, err := c.lockClient.Acquire(ctx, id)
	if err != nil {
		err := &er{
			err: errors.Wrap(err, "failed to acquire database lock"),
		}
		if errors.Is(err, dplock.ErrMongoDbClosing) {
			err.unavailable = true
		} else {
			err.conflict = true
		}
	}

	return l, err
}

func (c *Collection) unlock(ctx context.Context, lockID string) {
	c.lockClient.Unlock(ctx, lockID)
}
