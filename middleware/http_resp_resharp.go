package middleware

import (
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
)

type HttpRespReshapeMiddleware struct {
	BaseMiddleware
}

func NewRespReshapeMiddleware() *HttpRespReshapeMiddleware {
	return &HttpRespReshapeMiddleware{}
}

func (mw *HttpRespReshapeMiddleware) Priority() int {
	return 1000
}

func (mw *HttpRespReshapeMiddleware) Handle(ctx droplet.Context) error {
	code, message := 0, ""
	var d interface{}
	if err := mw.BaseMiddleware.Handle(ctx); err != nil {
		switch t := err.(type) {
		case *data.BaseError:
			code, message, d = t.Code, t.Message, t.Data
		default:
			code, message = data.ErrCodeInternal, err.Error()
		}
		var resp droplet.HttpResponse
		if r, ok := ctx.Output().(droplet.HttpResponse); ok {
			resp = r
		} else {
			resp = droplet.Option.ResponseNewFunc()
		}
		resp.Set(code, message, d)
		resp.SetReqID(ctx.GetString(KeyRequestID))
		ctx.SetOutput(resp)
		// response reshape is the last step, so we don't need return error
		return nil
	}

	switch ctx.Output().(type) {
	case droplet.HttpFileResponse:
	case droplet.HttpResponse:
		resp := ctx.Output().(droplet.HttpResponse)
		resp.SetReqID(ctx.GetString(KeyRequestID))
	default:
		resp := droplet.Option.ResponseNewFunc()
		resp.Set(code, message, ctx.Output())
		resp.SetReqID(ctx.GetString(KeyRequestID))
		ctx.SetOutput(resp)
	}

	return nil
}
