package logs

import "net/http"

func NewLoggingMiddleware(handler http.Handler) http.Handler {
	// AK TODO implement
	return handler
}

func NewTraceIdMiddleware(handler http.Handler) http.Handler {
	// AK TODO implement
	return nil
}
