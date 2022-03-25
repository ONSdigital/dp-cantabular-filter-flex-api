package api

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

const (
	DefaultLimit = 20
	DefaultOffset = 0
)

// getPaginationParams parses a URL and extracts limit/offset values.
func getPaginationParams(url *url.URL, maximumLimit int) (int, int, error) {
	query := url.Query()

	limit, err := getInt(query.Get("limit"), DefaultLimit)
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid parameter limit")
	}

	if limit > maximumLimit {
		return 0, 0, errors.Errorf("limit cannot be larger than %v", maximumLimit)
	}

	if limit < 0 {
		return 0, 0, errors.New("limit cannot be less than 0")
	}

	offset, err := getInt(query.Get("offset"), DefaultOffset)
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid parameter offset")
	}

	if offset < 0 {
		return 0, 0, errors.New("offset cannot be less than 0")
	}

	return limit, offset, nil
}

// getInt attempts to parse a passed string into a number, falling back to
// a default value if no string is provided.
func getInt(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot convert %s to int", value)
	}

	return parsed, nil
}
