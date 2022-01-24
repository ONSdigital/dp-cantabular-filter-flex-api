package mongodb

import (
	"fmt"
	"context"

	lock "github.com/ONSdigital/dp-mongodb/v3/dplock"
	"github.com/ONSdigital/dp-mongodb/v3/health"
	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const (
	filtersCollection = "filters"
)

// Client is the client responsible for querying mongodb 
type Client struct{
	conn   *mongo.MongoConnection
	health *health.CheckMongoClient
	lock   *lock.Lock
	cfg    Config
}

// NewClient returns a new mongodb Client
func NewClient(ctx context.Context, cfg Config) (*Client, error){
	c := Client{
		cfg: cfg,
	}
	var err error

	if c.conn, err = mongo.Open(&cfg.MongoDriverConfig); err != nil {
		return nil, fmt.Errorf("failed to open mongodb connection: %w", err)
	}

	collectionBuilder := map[health.Database][]health.Collection{
		health.Database(cfg.Database): {
			filtersCollection,
		},
	}

	c.health = health.NewClientWithCollections(c.conn, collectionBuilder)

	c.lock = lock.New(ctx, c.conn, filtersCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to create lock client: %w", err)
	}

	return &c, nil
}

// Close represents mongo session closing within the context deadline
func (c *Client) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (c *Client) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return c.health.Checker(ctx, state)
}