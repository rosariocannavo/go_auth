package models

import (
	"bytes"
	"net/http"
)

// ResponseRecorderWriter is a custom ResponseWriter that captures the response sended by the real server to original client(used for Proxy)
type ResponseRecorderWriter struct {
	http.ResponseWriter
	Body        *bytes.Buffer
	StatusCode  int
	StatusText  string
	wroteHeader bool
}

// NewResponseRecorderWriter creates a new ResponseRecorderWriter
func NewResponseRecorderWriter(w http.ResponseWriter) *ResponseRecorderWriter {
	return &ResponseRecorderWriter{
		ResponseWriter: w,
		Body:           bytes.NewBuffer(nil),
		StatusCode:     http.StatusOK, // Default status code
		StatusText:     http.StatusText(http.StatusOK),
	}
}

// WriteHeader captures the status code and message
func (rrw *ResponseRecorderWriter) WriteHeader(code int) {
	if !rrw.wroteHeader {
		rrw.StatusCode = code
		rrw.StatusText = http.StatusText(code)
		rrw.ResponseWriter.WriteHeader(code)
		rrw.wroteHeader = true
	}
}

// Write writes to the buffer and the original ResponseWriter
func (rrw *ResponseRecorderWriter) Write(b []byte) (int, error) {
	// Write to the buffer
	n, err := rrw.Body.Write(b)
	if err != nil {
		return n, err
	}

	// Write to the original ResponseWriter
	if !rrw.wroteHeader {
		rrw.WriteHeader(http.StatusOK)
	}

	_, _ = rrw.ResponseWriter.Write(b)
	return n, nil
}
