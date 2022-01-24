package responder

type dataLogger interface {
	LogData() map[string]interface{}
}

type coder interface {
	Code() int
}

type messager interface {
	Message() string
}
