package middlewares

import "github.com/shiningrush/droplet"

type InputMiddleWare struct {
	BaseMiddleWares
	input interface{}
}

func NewInputMiddleWare(input interface{}) *InputMiddleWare {
	return &InputMiddleWare{input: input}
}

func (mw *InputMiddleWare) Handle(context droplet.Context) error {
	context.SetInput(mw.input)
	return mw.next.Handle(context)
}
