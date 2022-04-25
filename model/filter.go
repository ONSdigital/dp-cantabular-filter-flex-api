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
	InstanceID        string              `bson:"instance_id"                  json:"instance_id"`
	Dataset           Dataset             `bson:"dataset"                      json:"dataset"`
	Published         bool                `bson:"published"                    json:"published"`
	DisclosureControl *DisclosureControl  `bson:"disclosure_control,omitempty" json:"disclosure_control,omitempty"`
	Type              string              `bson:"type"                         json:"type"`
	PopulationType    string              `bson:"population_type"              json:"population_type"`
	Links             FilterLinks         `bson:"links"                        json:"links"`
	Dimensions        []Dimension         `bson:"dimensions"                   json:"dimensions,omitempty"`
	UniqueTimestamp   primitive.Timestamp `bson:"unique_timestamp"             json:"-"`
	LastUpdated       time.Time           `bson:"last_updated"                 json:"-"`
	ETag              string              `bson:"etag"                         json:"-"`
}

type FilterLinks struct {
	Version    Link `bson:"version"    json:"version"`
	Self       Link `bson:"self"       json:"self"`
	Dimensions Link `bson:"dimensions" json:"dimensions"`
}

type FilterOutputLinks struct {
	Version         Link `bson:"version" json:"version"`
	Self            Link `bson:"self"    json:"self"`
	FilterBlueprint Link `json:"filter_blueprint"`
}

type SubmitFilterLinks struct {
	Dimensions Link `bson:"dimensions" json:"dimensions"`
}

type Link struct {
	HREF string `bson:"href"           json:"href"`
	ID   string `bson:"id,omitempty"   json:"id,omitempty"`
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
	Timestamp string `bson:"timestamp" json:"timestamp"`
	Name      string `bson:"name"      json:"name"`
}

type Dimension struct {
	Name       string   `bson:"name"          json:"name"`
	ID         string   `bson:"id"            json:"id"`
	Label      string   `bson:"label"         json:"label"`
	Options    []string `bson:"options"       json:"options"`
	IsAreaType bool     `bson:"is_area_type"  json:"is_area_type"`
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
		return errors.New(`"private" is empty in input`)
	}

	if len(strings.Trim(fi.Public, cutset)) == 0 {
		return errors.New(`"public" is empty in input`)
	}

	if len(strings.Trim(fi.Size, cutset)) == 0 {
		return errors.New(`"size" is empty in input`)
	}

	return nil
}
