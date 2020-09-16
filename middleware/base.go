package middleware

import (
	"github.com/shiningrush/droplet"
)

type BaseMiddleware struct {
	next droplet.Middleware
}

func (mw *BaseMiddleware) SetNext(next droplet.Middleware) {
	mw.next = next
}

func (mw *BaseMiddleware) Handle(ctx droplet.Context) error {
	return mw.next.Handle(ctx)
}
