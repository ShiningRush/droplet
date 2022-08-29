package droplet

import (
	"github.com/shiningrush/droplet/codec"
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"
)

type (
	Context = core.Context
	Handler = core.Handler
)

var (
	NewContext = core.NewContext
)

var (
	Option = GlobalOpt{
		HeaderKeyRequestID: "X-Request-ID",
		ResponseNewFunc: func() core.HttpResponse {
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
	ResponseNewFunc    func() core.HttpResponse
	Orchestrator       core.Orchestrator
	Codec              []codec.Interface
	TrafficLogOpt      *middleware.TrafficLogOpt
}
