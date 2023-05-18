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
		ResponseNewFunc: func() data.HttpResponse {
			return &data.Response{}
		},
		Codec: []codec.Interface{
			&codec.Json{},
			&codec.MultipartForm{},
			&codec.Empty{},
		},
		ErrSetting: ErrSetting{
			DefaultErrCode:  data.ErrCodeInternal,
			ValidateErrCode: data.ErrCodeValidate,
		},
	}
)

type GlobalOpt struct {
	HeaderKeyRequestID string
	ResponseNewFunc    func() data.HttpResponse
	Orchestrator       core.Orchestrator
	Codec              []codec.Interface
	TrafficLogOpt      *middleware.TrafficLogOpt
	ErrSetting         ErrSetting
}

type ErrSetting struct {
	DefaultErrCode  int
	ValidateErrCode int
}

func AppendCodec(code codec.Interface) {
	Option.Codec = append(Option.Codec, code)
}

func WithDefaultErrCode(code int) {
	Option.ErrSetting.DefaultErrCode = code
}

func WithValidateErrCode(code int) {
	Option.ErrSetting.ValidateErrCode = code
}
