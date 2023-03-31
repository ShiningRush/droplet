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
	handlerErr := mw.BaseMiddleware.Handle(ctx)

	var resp data.HttpResponse
	switch t := ctx.Output().(type) {
	case data.RawHttpResponse, data.HttpFileResponse:
		return nil
	case data.HttpResponse:
		resp = t
	default:
		// wrap result
		resp = mw.respNewFunc()
		resp.Set(0, "", ctx.Output())
		ctx.SetOutput(resp)
	}

	resp.SetReqID(ctx.GetString(KeyRequestID))
	if handlerErr != nil {
		switch t := handlerErr.(type) {
		case *data.BaseError:
			resp.Set(t.Code, t.Message, t.Data)
		default:
			resp.Set(data.ErrCodeInternal, handlerErr.Error(), nil)
		}
	}

	// response reshape is the last step, so we don't need to return error
	return nil
}
