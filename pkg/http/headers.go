package http

import "net/http"

// RequestOption is a function that mutates an outgoing request before it is sent.
type RequestOption func(*http.Request)

// WithHeaders returns a RequestOption that sets each key/value pair on the request.
func WithHeaders(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
}

// WithBearerToken returns a RequestOption that sets the Authorization header.
func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}
