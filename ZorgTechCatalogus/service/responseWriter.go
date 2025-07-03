package service

import "net/http"

type ServiceResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func NewServiceResponseWriter() *ServiceResponseWriter {
	return &ServiceResponseWriter{
		header: http.Header{},
	}
}

func (w *ServiceResponseWriter) Header() http.Header {
	return w.header
}

func (w *ServiceResponseWriter) Write(b []byte) (int, error) {
	if w.body == nil {
		w.body = b
	} else {
		w.body = append(w.body, b...)
	}
	return 0, nil
}

func (w *ServiceResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

var okFn = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
