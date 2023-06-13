package mongodb

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	lock "github.com/ONSdigital/dp-mongodb/v3/dplock"
	"github.com/ONSdigital/dp-mongodb/v3/health"
	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"github.com/pkg/errors"
)

// Client is the client responsible for querying mongodb
type Client struct {
	conn        *mongo.MongoConnection
	health      *health.CheckMongoClient
	lock        *lock.Lock
	cfg         Config
	collections *Collections
	generate    generator
}

// NewClient returns a new mongodb Client
func NewClient(ctx context.Context, g generator, cfg Config) (*Client, error) {
	c := Client{
		cfg:      cfg,
		generate: g,
	}
	var err error

	if c.conn, err = mongo.Open(&cfg.MongoDriverConfig); err != nil {
		return nil, errors.Wrap(err, "failed to open mongodb connection: %w")
	}

	collectionBuilder := map[health.Database][]health.Collection{
		health.Database(cfg.Database): {
			health.Collection(cfg.FiltersCollection),
			health.Collection(cfg.FilterOutputsCollection),
		},
	}

	c.health = health.NewClientWithCollections(c.conn, collectionBuilder)

	c.collections = &Collections{
		filters: &Collection{
			name:       cfg.FiltersCollection,
			lockClient: lock.New(ctx, c.conn, cfg.FiltersCollection),
		},
		filterOutputs: &Collection{
			name:       cfg.FilterOutputsCollection,
			lockClient: lock.New(ctx, c.conn, cfg.FilterOutputsCollection),
		},
	}

	return &c, nil
}

// Conn returns the underlying mongodb connection.
func (c *Client) Conn() *mongo.MongoConnection {
	return c.conn
}

// Close represents mongo session closing within the context deadline
func (c *Client) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (c *Client) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return c.health.Checker(ctx, state)
}
