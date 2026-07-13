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
	KeyHttpRequest        = "HttpRequest"
	KeyTraceID            = "ofa-pass-trace-id"
	KeyRequestID          = "ofa-direct-request-id"
	KeyOperator           = "ofa-pass-operator"
	KeyTenantID           = "ofa-pass-tenant-id"
	KeyAppID              = "ofa-pass-app-id"
	KeyLocale             = "ofa-pass-locale"
	KeyRemainingTimeoutMS = "ofa-direct-remaining-timeout-ms"
)
