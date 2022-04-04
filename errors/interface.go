package errors

type errNotFound interface {
	NotFound() bool
}

type errConflict interface {
	Conflict() bool
}

type errUnavailable interface {
	Unavailable() bool
}

type errForbidden interface {
	Forbidden() bool
}
