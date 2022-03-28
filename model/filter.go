package model

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Completed = "completed"
	Submitted = "submitted"
)

// Filter holds details for a user filter journey
type Filter struct {
	ID                string              `bson:"filter_id"                    json:"filter_id"`
	Links             Links               `bson:"links"                        json:"links"`
	UniqueTimestamp   primitive.Timestamp `bson:"unique_timestamp"             json:"-"`
	LastUpdated       time.Time           `bson:"last_updated"                 json:"-"`
	ETag              string              `bson:"etag"                         json:"-"`
	InstanceID        string              `bson:"instance_id"                  json:"instance_id"`
	Dimensions        []Dimension         `bson:"dimensions"                   json:"-"`
	Dataset           Dataset             `bson:"dataset"                      json:"dataset"`
	Published         bool                `bson:"published"                    json:"published"`
	DisclosureControl *DisclosureControl  `bson:"disclosure_control,omitempty" json:"disclosure_control,omitempty"`
	Type              string              `bson:"type"                         json:"type"`
	PopulationType    string              `bson:"population_type"              json:"population_type"`
}

// PutFilter holds details for PUT filter response
type PutFilter struct {
	Events         []Event `bson:"events"                       json:"events"`
	Dataset        Dataset `bson:"dataset"                      json:"dataset"`
	PopulationType string  `bson:"population_type"              json:"population_type"`
}

type JobState struct {
	InstanceID       string  `json:"instance_id"`
	DimensionListUrl string  `json:"dimension_list_url"`
	FilterID         string  `json:"filter_id"`
	Events           []Event `json:"events"`
}

type Links struct {
	Version Link `bson:"version" json:"version"`
	Self    Link `bson:"self"    json:"self"`
}

type FilterOutputLinks struct {
	Links
	FilterBlueprint Link `json:"filter_blueprint"`
}

type Link struct {
	HREF string `bson:"href"           json:"href"`
	ID   string `bson:"id,omitempty"   json:"id,omitempty"`
}

type FilterOutput struct {
	ID                string              `bson:"id,omitempty"                 json:"id,omitempty"`
	FilterID          string              `bson:"filter_id"                    json:"filter_id"`
	InstanceID        string              `bson:"instance_id"                  json:"instance_id"`
	Dataset           Dataset             `bson:"dataset"                      json:"dataset"`
	Published         bool                `bson:"published"                    json:"published"`
	State             string              `bson:"state,omitempty"              json:"state,omitempty"`
	Downloads         Downloads           `bson:"downloads"                    json:"downloads"`
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

type Event struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Name      string    `bson:"name"      json:"name"`
}

type Dimension struct {
	Name         string   `bson:"name"          json:"name"`
	Options      []string `bson:"options"       json:"options"`
	DimensionURL string   `bson:"dimension_url" json:"dimension_url,omitempty"`
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

func (fi *FileInfo) IsNotFullyPopulated() error {
	cutset := " "

	if len(strings.Trim(fi.HREF, cutset)) == 0 {
		return errors.New(`"HREF" is empty in input`)
	}

	if len(strings.Trim(fi.Private, cutset)) == 0 {
		return errors.New(`"Private" is empty in input`)
	}

	if len(strings.Trim(fi.Public, cutset)) == 0 {
		return errors.New(`"Public" is empty in input`)
	}

	if len(strings.Trim(fi.Size, cutset)) == 0 {
		return errors.New(`"Size" is empty in input`)
	}

	return nil
}
