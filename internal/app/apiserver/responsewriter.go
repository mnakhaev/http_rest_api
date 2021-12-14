package apiserver

import "net/http"

type responseWriter struct {
	// using anonymous field. No need to realize all the methods from response writer - they will be already available
	http.ResponseWriter
	code int
}

// Redefine WriteHeader method because it writes status code in original http package
// It's like a Python decorator. Keep current behaviour of WriteHeader, but just add status code
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
