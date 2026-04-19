package core_http_response

import "net/http"

var (
	statusCodeUnintialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     statusCodeUnintialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// Write ensures that if the handler writes a body without explicitly
// calling WriteHeader, we record the default status code (200) so
// tracing middleware can read it safely.
// func (rw *ResponseWriter) Write(b []byte) (int, error) {
// 	if rw.statusCode == statusCodeUnintialized {
// 		rw.statusCode = http.StatusOK
// 	}
// 	return rw.ResponseWriter.Write(b)
// }

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {

	if rw.statusCode == statusCodeUnintialized {
		panic("status code not initialized")
	}

	return rw.statusCode
}
