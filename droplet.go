package droplet

import "github.com/shiningrush/droplet/data"

var (
	Option = GlobalOpt{
		HeaderKeyRequestID: "X-Request-ID",
		ResponseNewFunc: func() HttpResponse {
			return &data.Response{}
		},
	}
)

type GlobalOpt struct {
	HeaderKeyRequestID string
	ResponseNewFunc    func() HttpResponse
	Orchestrator       Orchestrator
}

type HttpResponse interface {
	Set(code int, msg, reqId string, data interface{})
}

type HttpFileResponse interface {
	Get() (name, contentType string, content []byte)
}
