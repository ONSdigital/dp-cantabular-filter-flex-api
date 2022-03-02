package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Filter holds details for a user filter journey
type Filter struct {
	ID                string               `bson:"filter_id"                    json:"filter_id"`
	Links             Links                `bson:"links"                        json:"links"`
	FilterOutput      *FilterOutput        `bson:"filter_output,omitempty"      json:"filter_output,omitempty"`
	Events            []Event              `bson:"events"                       json:"events"`
	UniqueTimestamp   primitive.Timestamp  `bson:"unique_timestamp"             json:"-"`
	LastUpdated       time.Time            `bson:"last_updated"                 json:"-"`
	ETag              string               `bson:"etag"                         json:"-"`
	InstanceID        string               `bson:"instance_id"                  json:"instance_id"`
	Dimensions        []Dimension          `bson:"dimensions"                   json:"dimensions"`
	Dataset           Dataset              `bson:"dataset"                      json:"dataset"`
	Published         bool                 `bson:"published"                    json:"published"`
	DisclosureControl *DisclosureControl   `bson:"disclosure_control,omitempty" json:"disclosure_control,omitempty"`
	Type              string               `bson:"type"                         json:"type"`
	PopulationType    string               `bson:"population_type"              json:"population_type"`
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
	CSV  *FileInfo `bson:"csv,omitempty"  json:"csv,omitempty"`
	CSVW *FileInfo `bson:"csvw,omitempty" json:"csvw,omitempty"`
	TXT  *FileInfo `bson:"txt,omitempty"  json:"txt,omitempty"`
	XLS  *FileInfo `bson:"xls,omitempty"  json:"xls,omitempty"`
}

type FileInfo struct {
	HREF    string `bson:"href"    json:"href"`
	Size    string `bson:"size"    json:"size"`
	Public  string `bson:"public"  json:"public"`
	Private string `bson:"private" json:"private"`
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
