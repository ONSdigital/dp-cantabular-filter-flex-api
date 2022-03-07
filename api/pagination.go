package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

const DefaultLimit = 20
const DefaultOffset = 0

// getPaginationParams parses a URL and extracts limit/offset values.
func getPaginationParams(url *url.URL, maximumLimit int, logData log.Data) (limit, offset int, err error) {
	query := url.Query()

	limit, err = getInt(query.Get("limit"), DefaultLimit)
	if err != nil {
		return 0, 0, &Error{err: err, message: "invalid parameter limit", logData: logData}
	}

	if limit > maximumLimit {
		msg := fmt.Sprintf("limit cannot be larger than %v", maximumLimit)
		return 0, 0, &Error{err: errors.New(msg), message: msg, logData: logData}
	}

	if limit < 0 {
		msg := "limit cannot be less than 0"
		return 0, 0, &Error{err: errors.New(msg), message: msg, logData: logData}
	}

	offset, err = getInt(query.Get("offset"), DefaultOffset)
	if err != nil {
		return 0, 0, &Error{err: err, message: "invalid parameter offset", logData: logData}
	}

	if offset < 0 {
		msg := "offset cannot be less than 0"
		return 0, 0, &Error{err: errors.New(msg), message: msg, logData: logData}
	}

	return
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
