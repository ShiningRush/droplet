package wrapper

import "net/http"

// NewResponseWriter wrap a http.ResponseWriter to ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *StdWriterWrapper {
	return &StdWriterWrapper{stdWriter: w}
}

type StdWriterWrapper struct {
	stdWriter http.ResponseWriter
}

func (r *StdWriterWrapper) SetHeader(key, val string) {
	r.stdWriter.Header().Set(key, val)
}

func (r *StdWriterWrapper) GetHeader(key string) string {
	return r.stdWriter.Header().Get(key)
}

func (r *StdWriterWrapper) GetHeaderValues(key string) []string {
	return r.stdWriter.Header().Values(key)
}

func (r *StdWriterWrapper) DelHeader(key string) {
	r.stdWriter.Header().Del(key)
}

func (r *StdWriterWrapper) Write(bs []byte) (int, error) {
	return r.stdWriter.Write(bs)
}

func (r *StdWriterWrapper) WriteHeader(statusCode int) {
	r.stdWriter.WriteHeader(statusCode)
}

func (r *StdWriterWrapper) StdHttpWriter() http.ResponseWriter {
	return r.stdWriter
}
