package middlewares

import "github.com/shiningrush/droplet"

type InputMiddleware struct {
	BaseMiddleware
	input interface{}
}

func NewInputMiddleWare(input interface{}) *InputMiddleware {
	return &InputMiddleware{input: input}
}

func (mw *InputMiddleware) Handle(context droplet.Context) error {
	context.SetInput(mw.input)
	return mw.next.Handle(context)
}
