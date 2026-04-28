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

func (rw *ResponseWriter) GetStatusCode() int {

	if rw.statusCode == statusCodeUnintialized {
		// panic("status code not initialized")
		return http.StatusOK
	}

	return rw.statusCode
}
