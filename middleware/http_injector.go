package middleware

import (
	"github.com/shiningrush/droplet"
	"net/http"
)

type HttpInfoInjectorMiddleware struct {
	BaseMiddleware
	opt HttpInfoInjectorOption
}

type HttpInfoInjectorOption struct {
	ReqFunc func() *http.Request
}

func NewHttpInfoInjectorMiddleware(opt HttpInfoInjectorOption) *HttpInfoInjectorMiddleware {
	return &HttpInfoInjectorMiddleware{opt: opt}
}

func (mw *HttpInfoInjectorMiddleware) Handle(ctx droplet.Context) error {
	req := mw.opt.ReqFunc()

	ctx.Set(KeyHttpRequest, req)
	ctx.Set(KeyRequestID, req.Header.Get(droplet.Option.HeaderKeyRequestID))
	ctx.SetPath(req.URL.Path)

	return mw.BaseMiddleware.Handle(ctx)
}
