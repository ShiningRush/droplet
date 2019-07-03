package middlewares

import "github.com/shiningrush/droplet"

type BaseMiddleWares struct {
	next droplet.MiddleWare
}

func(mw *BaseMiddleWares) SetNext(next droplet.MiddleWare){
	mw.next = next
}

func(mw *BaseMiddleWares) Handle(ctx droplet.Context) error{
	return mw.next.Handle(ctx)
}
