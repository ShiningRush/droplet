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
	ctx.Set(KeyHttpRequest, mw.opt.ReqFunc())
	ctx.Set(KeyRequestID, mw.opt.ReqFunc().Header.Get(droplet.Option.HeaderKeyRequestID))
	ctx.SetPath(mw.opt.ReqFunc().URL.Path)

	return mw.BaseMiddleware.Handle(ctx)
}
