package droplet

type MiddleWare interface {
	SetNext(next MiddleWare)
	Handle(context Context) error
}

type handlerMiddleWare struct {
	handler Handler
	next    MiddleWare
}

func (m *handlerMiddleWare) SetNext(next MiddleWare) {
	m.next = next
}

func (m *handlerMiddleWare) Handle(context Context) error {
	rs, err := m.handler()
	context.SetOutput(rs)
	return err
}

func NewHandlerMiddleWare(handler Handler) MiddleWare {
	return &handlerMiddleWare{handler: handler}
}
