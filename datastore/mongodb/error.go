package mongodb

// er is the packages error type
type er struct {
	err         error
	logData     map[string]interface{}
	notFound    bool
	conflict    bool
	unavailable bool
}

// Error satisfies the standard library Go error interface
func (e *er) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// Unwrap implements the standard library Go unwrapper interface
func (e *er) Unwrap() error {
	return e.err
}

// LogData satisfies the dataLogger interface which is used to recover
// log data from an error
func (e *er) LogData() map[string]interface{} {
	return e.logData
}

// NotFound satisfies the errNotFound interface and allows other packages
// to recall metadata about the error thrown
func (e *er) NotFound() bool {
	return e.notFound
}

// Conflict satisfies the errConflict interface and allows other packages
// to recall metadata about the error thrown
func (e *er) Conflict() bool {
	return e.conflict
}

// Unavailable satisfies the errUnavailable interface and allows other packages
// to recall metadata about the error thrown
func (e *er) Unavailable() bool {
	return e.unavailable
}
