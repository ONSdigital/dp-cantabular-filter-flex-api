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

type Downloads struct {
	CSV  *FileInfo `bson:"csv,omitempty"   json:"csv,omitempty"`
	CSVW *FileInfo `bson:"csvw,omitempty"  json:"csvw,omitempty"`
	TXT  *FileInfo `bson:"txt,omitempty"   json:"txt,omitempty"`
	XLS  *FileInfo `bson:"xls,omitempty"   json:"xls,omitempty"`
}

type FileInfo struct {
	HREF    string `bson:"href"    json:"href"`
	Size    string `bson:"size"    json:"size"`
	Public  string `bson:"public"  json:"public"`
	Private string `bson:"private" json:"private"`
	Skipped bool   `bson:"skipped" json:"skipped"`
}
