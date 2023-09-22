package middleware

import (
	"net/http"

	"github.com/shiningrush/droplet/core"
)

type HttpInfoInjectorMiddleware struct {
	BaseMiddleware
	opt HttpInfoInjectorOption
}

type HttpInfoInjectorOption struct {
	ReqFunc            func() *http.Request
	HeaderKeyRequestID string
}

func NewHttpInfoInjectorMiddleware(opt HttpInfoInjectorOption) *HttpInfoInjectorMiddleware {
	return &HttpInfoInjectorMiddleware{opt: opt}
}

func (mw *HttpInfoInjectorMiddleware) Handle(ctx core.Context) error {
	req := mw.opt.ReqFunc()

	ctx.Set(KeyHttpRequest, req)
	if ctx.Get(KeyRequestID) == nil {
		ctx.Set(KeyRequestID, req.Header.Get(mw.opt.HeaderKeyRequestID))
	}
	ctx.SetPath(req.URL.Path)

	return mw.BaseMiddleware.Handle(ctx)
}
