package service

import (
	"net/http"
)

// NotImplemented replies to the request with an HTTP 501 Not Implemented error.
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "开发中...", http.StatusNotImplemented)
}

// NotImplementedHandler returns a simple request handler
// that replies to each request with a ``开发中...'' reply.
func NotImplementedHandler() http.Handler {
	return http.HandlerFunc(NotImplemented)
}
