package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

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
		reqId := req.Header.Get(mw.opt.HeaderKeyRequestID)
		if reqId == "" {
			reqId = fmt.Sprintf("%s%v", time.Now().Format("20060102150405"), rand.Intn(100000))
		}
		ctx.Set(KeyRequestID, reqId)
	}
	ctx.SetPath(req.URL.Path)

	return mw.BaseMiddleware.Handle(ctx)
}
