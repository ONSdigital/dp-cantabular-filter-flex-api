package responder

import (
	"errors"
	"strconv"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ONSdigital/log.go/v2/log"
)

// unwrapStatusCode is a callback function that allows you to extract
// a status code from an error. If the status code returned is 0,
// statusCode will attempt to recursively unwrap the error until a
// non-zero code is returned. If no more code is embedded it will
// return status 500 as default.
func unwrapStatusCode(err error) int {
	if code := statusCode(err); code != 0 {
		return code
	}

	for errors.Unwrap(err) != nil {
		if code := statusCode(err); code != 0 {
			return code
		}
		err = errors.Unwrap(err)
	}

	return http.StatusInternalServerError
}

// statusCode attempts to extract a status code from an error,
// or returns 0 if not found
func statusCode(err error) int {
	var cerr coder

	if errors.As(err, &cerr) {
		return cerr.Code()
	}

	return 0
}

// logData returns logData for an error if there is any. This is used
// to extract log.Data embedded in an error if it implements the dataLogger
// interface
func logData(err error) log.Data {
	var lderr dataLogger
	if errors.As(err, &lderr) {
		return lderr.LogData()
	}

	return nil
}

// unwrapLogData recursively unwraps logData from an error. This allows an
// error to be wrapped with log.Data at each level of the call stack, and
// then extracted and combined here as a single log.Data entry. This allows
// us to log errors only once but maintain the context provided by log.Data
// at each level.
func unwrapLogData(err error) log.Data {
	var data []log.Data

	for err != nil && errors.Unwrap(err) != nil {
		if lderr, ok := err.(dataLogger); ok {
			if d := lderr.LogData(); d != nil {
				data = append(data, d)
			}
		}

		err = errors.Unwrap(err)
	}

	if len(data) == 0{
		return nil
	}

	// flatten []log.Data into single log.Data with slice
	// entries for duplicate keyed entries, but not for duplicate
	// key-value pairs
	logData := log.Data{}
	for _, d := range data {
		for k, v := range d {
			if val, ok := logData[k]; ok {
				if !reflect.DeepEqual(val, v) {
					if s, ok := val.([]interface{}); ok {
						s = append(s, v)
						logData[k] = s
					} else {
						logData[k] = []interface{}{val, v}
					}
				}
			} else {
				logData[k] = v
			}
		}
	}

	return logData
}

// errorResponse extracts a specified error response to be returned
// to the caller if present, otherwise returns the default error
// string
func errorMessage(err error) string {
	var rerr messager
	if errors.As(err, &rerr) {
		if resp := rerr.Message(); resp != "" {
			return resp
		}
	}

	return err.Error()
}

// stackTrace recursively unwraps the error looking for the deepest
// level at which the error was wrapped with a stack trace from
// github.com/pkg/errors (or conforms to the StackTracer interface)
// and returns the slice of stack frames. These can are of type
// log.go/EventStackTrace so can be used directly with log.Go's
// available API to preserve the correct error logging format
func stackTrace(err error) []log.EventStackTrace{
	var serr stacktracer
	var resp []log.EventStackTrace

	for errors.Unwrap(err) != nil && errors.As(err, &serr) {
		st := serr.StackTrace()
		resp = make([]log.EventStackTrace, 0)
		for _, f := range st{
			line, _ := strconv.Atoi(fmt.Sprintf("%d",  f))
			resp = append(resp, log.EventStackTrace{
				File:     fmt.Sprintf("%+s", f),
				Function: fmt.Sprintf("%n",  f),
				Line:     line,
			})
		}

		err = errors.Unwrap(err)
	}

	return resp
}
