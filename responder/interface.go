package responder

import (
	"github.com/pkg/errors"
)

type dataLogger interface {
	LogData() map[string]interface{}
}

type coder interface {
	Code() int
}

type messager interface {
	Message() string
}

type stacktracer interface {
	StackTrace() errors.StackTrace
}
