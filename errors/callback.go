// package errors is a temporary package to hold the interfaces
// and callback for the service error types. The idea is for these
// to be moved out to dp-net/errors but am including here until
// those addtions are approved to prevent being held up
package errors

import (
	"errors"
)

func NotFound(err error) bool {
	var e errNotFound

	if errors.As(err, &e) {
		return e.NotFound()
	}

	return false
}

func Conflict(err error) bool {
	var e errConflict

	if errors.As(err, &e) {
		return e.Conflict()
	}

	return false
}

func Unavailable(err error) bool {
	var e errUnavailable

	if errors.As(err, &e) {
		return e.Unavailable()
	}

	return false
}

func Forbidden(err error) bool {
	var e errForbidden

	if errors.As(err, &e) {
		return e.Forbidden()
	}
	return false
}

func BadRequest(err error) bool {
	var e errBadRequest

	if errors.As(err, &e) {
		return e.BadRequest()
	}

	return false
}
