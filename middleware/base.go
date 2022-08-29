package middleware

import (
	"github.com/shiningrush/droplet/core"
)

type BaseMiddleware struct {
	next core.Middleware
}

func (mw *BaseMiddleware) SetNext(next core.Middleware) {
	mw.next = next
}

func (mw *BaseMiddleware) Handle(ctx core.Context) error {
	return mw.next.Handle(ctx)
}

const (
	KeyHttpRequest = "HttpRequest"
	KeyRequestID   = "RequestID"
)
