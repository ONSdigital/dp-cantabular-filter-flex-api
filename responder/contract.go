package responder

// errorResponse is the generic ONS error response for HTTP errors
type errorResponse struct {
	Errors []string `json:"errors"`
}
