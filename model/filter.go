package model

import (
	"time"

	"github.com/google/uuid"
)

// Filter holds details for a user filter journey
type Filter struct {
	ID                uuid.UUID          `bson:"filter_id"                    json:"filter_id"`
	Links             Links              `bson:"links"                        json:"links"`
	FilterOutput      *FilterOutput      `bson:"filter_output,omitempty"      json:"filter_output,omitempty"`
	Events            []Event            `bson:"events"                       json:"events"`
	UniqueTimestamp   time.Time          `bson:"unique_timestamp"             json:"unique_timestamp"`
	LastUpdated       time.Time          `bson:"last_updated"                 json:"last_updated"`
	ETag              string             `bson:"etag"                         json:"etag"`
	InstanceID        uuid.UUID          `bson:"instance_id"                  json:"instance_id"`
	Dimensions        []Dimension        `bson:"dimensions"                   json:"dimensions"`
	Dataset           Dataset            `bson:"dataset"                      json:"dataset"`
	Published         bool               `bson:"published"                    json:"published"`
	DisclosureControl *DisclosureControl `bson:"disclosure_control,omitempty" json:"disclosure_control,omitempty"`
	Type              string             `bson:"type"                         json:"type"`
	PopulationType    string             `bson:"population_type"              json:"population_type"`
}

type Links struct {
	Version Link `bson:"version" json:"version"`
	Self    Link `bson:"self"    json:"self"`
}

type Link struct {
	HREF string `bson:"href"           json:"href"`
	ID   string `bson:"id,omitempty"   json:"id,omitempty"`
}

type FilterOutput struct {
	CSV  *FileInfo `bson:"csv,omitempty"  json:"csv,omitempty" validate:"required"`
	CSVW *FileInfo `bson:"csvw,omitempty" json:"csvw,omitempty" validate:"required"`
	TXT  *FileInfo `bson:"txt,omitempty"  json:"txt,omitempty" validate:"required"`
	XLS  *FileInfo `bson:"xls,omitempty"  json:"xls,omitempty" validate:"required"`
}

type FileInfo struct {
	HREF    string `bson:"href"    json:"href" validate:"required"`
	Size    string `bson:"size"    json:"size" validate:"required"`
	Public  string `bson:"public"  json:"public" validate:"required"`
	Private string `bson:"private" json:"private" validate:"required"`
	Skipped bool   `bson:"skipped" json:"skipped"`
}

type Event struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Name      string    `bson:"name"      json:"name"`
}

type Dimension struct {
	Name         string   `bson:"name"          json:"name"`
	Options      []string `bson:"options"       json:"options"`
	DimensionURL string   `bson:"dimension_url" json:"dimension_url"`
	IsAreaType   bool     `bson:"is_area_type"  json:"is_area_type"`
}

type Dataset struct {
	ID      string `bson:"id"      json:"id"`
	Edition string `bson:"edition" json:"edition"`
	Version int    `bson:"version" json:"version"`
}

type DisclosureControl struct {
	Status         string         `bson:"status"          json:"status"`
	Dimension      string         `bson:"dimension"       json:"dimension"`
	BlockedOptions BlockedOptions `bson:"blocked_options" json:"blocked_options"`
}

type BlockedOptions struct {
	BlockedOptions []string `bson:"blocked_options" json:"blocked_options"`
	BlockedCount   int      `bson:"blocked_count"   json:"blocked_count"`
}
