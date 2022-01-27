package mongodb

import(
	"context"
	"github.com/ONSdigital/dp-mongodb/v3/dplock"
)

// Collections holds information about the mongodb collections
// relevant to this service
type Collections struct{
	filters       *Collection
	filterOutputs *Collection
}

// Collection represents a collection in mongodb. Holds
// The name(s), lock client and any other relevant information.
type Collection struct{
	name       string
	lockClient *dplock.Lock
}

func (c *Collection) lock(ctx context.Context, id string) (string, error){
	return c.lockClient.Acquire(ctx, id)
}

func (c *Collection) unlock(ctx context.Context, lockID string){
	c.lockClient.Unlock(ctx, lockID)
}
