package core

import (
	"net/http"

	"github.com/shiningrush/droplet/data"
)

type HttpResponse interface {
	Set(code int, msg string, data interface{})
	SetReqID(reqId string)
}

type HttpFileResponse interface {
	Get() *data.FileResponse
}

type SpecCodeHttpResponse interface {
	GetStatusCode() int
	HttpResponse
}

type ResponseWriter interface {
	SetHeader(key, val string)
	GetHeader(key string) string
	GetHeaderValues(key string) []string
	DelHeader(key string)

	Write([]byte) (int, error)
	WriteHeader(statusCode int)

	// StdHttpWriter return the http.ResponseWriter, if wrapped framework is not compatible(such as fasthttp)
	// it will return nil
	StdHttpWriter() http.ResponseWriter
}

type RawHttpResponse interface {
	WriteRawResponse(writer ResponseWriter) error
}
