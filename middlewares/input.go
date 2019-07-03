package middlewares

import "github.com/shiningrush/droplet"

type InputMiddleWare struct {
	next  droplet.MiddleWare
	input interface{}
}

func NewInputMiddleWare(input interface{}) *InputMiddleWare {
	return &InputMiddleWare{input: input}
}

func (mw *InputMiddleWare) SetNext(next droplet.MiddleWare) {
	mw.next = next
}

func (mw *InputMiddleWare) Handle(context droplet.Context) error {
	context.SetInput(mw.input)
	return mw.next.Handle(context)
}
