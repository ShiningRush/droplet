package droplet

import (
	"net/http"

	"github.com/shiningrush/droplet/codec"
	"github.com/shiningrush/droplet/data"
)

var (
	Option = GlobalOpt{
		HeaderKeyRequestID: "X-Request-ID",
		ResponseNewFunc: func() HttpResponse {
			return &data.Response{}
		},
		Codec: []codec.Interface{
			&codec.Json{},
			&codec.MultipartForm{},
			&codec.Empty{},
		},
	}
)

type GlobalOpt struct {
	HeaderKeyRequestID string
	ResponseNewFunc    func() HttpResponse
	Orchestrator       Orchestrator
	Codec              []codec.Interface
}

type HttpResponse interface {
	Set(code int, msg string, data interface{})
	SetReqID(reqId string)
}

type HttpFileResponse interface {
	Get() (name, contentType string, content []byte)
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
