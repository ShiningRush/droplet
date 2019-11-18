package droplet

type Handler func(ctx Context) (interface{}, error)

type Pipe interface {
	Add(mw Middleware) Pipe
	AddIf(mw Middleware, predicate bool) Pipe
	Run(handler Handler) (interface{}, error)
}

type BasePipe struct {
	mws   []Middleware
}

func NewPipe() *BasePipe {
	return &BasePipe{}
}

func (p *BasePipe) Add(mw Middleware) Pipe {
	p.mws = append(p.mws, mw)
	return p
}

func (p *BasePipe) AddIf(mw Middleware,predicate bool) Pipe {
	if predicate {
		p.Add(mw)
	}
	return p
}

func (p *BasePipe) AddRange(mws []Middleware) Pipe {
	for _, mw := range mws {
		p.Add(mw)
	}
	return p
}

func (p *BasePipe) Run(handler Handler) (interface{}, error) {
	initCtx := NewContext()

	handlerMw := NewHandlerMiddleware(handler)
	for i, mw := range p.mws {
		if i < len(p.mws)-1 {
			mw.SetNext(p.mws[i+1])
			continue
		}

		mw.SetNext(handlerMw)
	}

	if len(p.mws) == 0 {
		err := handlerMw.Handle(initCtx)
		return initCtx, err
	}

	err := p.mws[0].Handle(initCtx)
	return initCtx.output, err
}
