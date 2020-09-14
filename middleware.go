package droplet

type Middleware interface {
	Priority() int
	SetNext(next Middleware)
	Handle(ctx Context) error
}

type handlerMiddleware struct {
	handler Handler
	next    Middleware
}

func (m *handlerMiddleware) Priority() int {
	return -1
}

func (m *handlerMiddleware) SetNext(next Middleware) {
	m.next = next
}

func (m *handlerMiddleware) Handle(ctx Context) error {
	rs, err := m.handler(ctx)
	ctx.SetOutput(rs)
	return err
}

func NewHandlerMiddleware(handler Handler) Middleware {
	return &handlerMiddleware{handler: handler}
}
