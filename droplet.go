package droplet

import (
	"github.com/shiningrush/droplet/codec"
	"github.com/shiningrush/droplet/data"
	"net/http"
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

type RawHttpResponse interface {
	WriteRawResponse(http.ResponseWriter) error
}
