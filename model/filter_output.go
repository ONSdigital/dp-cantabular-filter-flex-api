package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterOutput struct {
	ID                string              `bson:"id,omitempty"                 json:"id,omitempty"`
	FilterID          string              `bson:"filter_id"                    json:"filter_id"`
	InstanceID        string              `bson:"instance_id"                  json:"instance_id"`
	Dataset           Dataset             `bson:"dataset"                      json:"dataset"`
	Published         bool                `bson:"published"                    json:"published"`
	State             string              `bson:"state,omitempty"              json:"state,omitempty"`
	Downloads         Downloads           `bson:"downloads,omitempty"          json:"downloads,omitempty"`
	Events            []Event             `bson:"events"                       json:"events"`
	Type              string              `bson:"type"                         json:"type"`
	PopulationType    string              `bson:"population_type"              json:"population_type"`
	DisclosureControl *DisclosureControl  `bson:"disclosure_control,omitempty" json:"disclosure_control,omitempty"`
	Links             FilterOutputLinks   `bson:"links"                        json:"links"`
	Dimensions        []Dimension         `bson:"dimensions"                   json:"dimensions"`
	UniqueTimestamp   primitive.Timestamp `bson:"unique_timestamp"             json:"-"`
	LastUpdated       time.Time           `bson:"last_updated"                 json:"-"`
	ETag              string              `bson:"etag"                         json:"-"`
}
