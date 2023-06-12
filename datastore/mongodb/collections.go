package mongodb

import (
	"github.com/ONSdigital/dp-mongodb/v3/dplock"
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
