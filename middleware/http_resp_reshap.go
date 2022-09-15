package middleware

import (
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
)

type HttpRespReshapeMiddleware struct {
	BaseMiddleware

	respNewFunc func() data.HttpResponse
}

func NewRespReshapeMiddleware(respNewFunc func() data.HttpResponse) *HttpRespReshapeMiddleware {
	return &HttpRespReshapeMiddleware{respNewFunc: respNewFunc}
}

func (mw *HttpRespReshapeMiddleware) Handle(ctx core.Context) error {
	code, message := 0, ""
	var d interface{}
	if err := mw.BaseMiddleware.Handle(ctx); err != nil {
		switch t := err.(type) {
		case *data.BaseError:
			code, message, d = t.Code, t.Message, t.Data
		default:
			code, message = data.ErrCodeInternal, err.Error()
		}
		var resp data.HttpResponse
		if r, ok := ctx.Output().(data.HttpResponse); ok {
			resp = r
		} else {
			resp = mw.respNewFunc()
		}
		resp.Set(code, message, d)
		resp.SetReqID(ctx.GetString(KeyRequestID))
		ctx.SetOutput(resp)
		// response reshape is the last step, so we don't need to return error
		return nil
	}

	switch ctx.Output().(type) {
	case data.RawHttpResponse, data.HttpFileResponse:
	case data.HttpResponse:
		resp := ctx.Output().(data.HttpResponse)
		resp.SetReqID(ctx.GetString(KeyRequestID))
	default:
		resp := mw.respNewFunc()
		resp.Set(code, message, ctx.Output())
		resp.SetReqID(ctx.GetString(KeyRequestID))
		ctx.SetOutput(resp)
	}

	return nil
}
