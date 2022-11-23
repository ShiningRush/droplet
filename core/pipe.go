package core

type Orchestrator func(mws []Middleware) []Middleware

type Handler func(ctx Context) (interface{}, error)

type BasePipe struct {
	globalOrchestrator Orchestrator
	mwOrchestrator     Orchestrator
	mws                []Middleware
}

func NewPipe(gOrc Orchestrator) *BasePipe {
	return &BasePipe{
		globalOrchestrator: gOrc,
	}
}

func (p *BasePipe) SetOrchestrator(o Orchestrator) *BasePipe {
	p.mwOrchestrator = o
	return p
}

func (p *BasePipe) Add(mw Middleware) *BasePipe {
	p.mws = append(p.mws, mw)
	return p
}

func (p *BasePipe) AddIf(mw Middleware, predicate bool) *BasePipe {
	if predicate {
		p.Add(mw)
	}
	return p
}

func (p *BasePipe) AddRange(mws []Middleware) *BasePipe {
	for _, mw := range mws {
		p.Add(mw)
	}
	return p
}

type RunOpt struct {
	InitContext Context
}

type SetRunOpt func(opt *RunOpt)

func WithRunContext(ctx Context) SetRunOpt {
	return func(opt *RunOpt) {
		opt.InitContext = ctx
	}
}

func (p *BasePipe) Run(handler Handler, opts ...SetRunOpt) (interface{}, error) {
	opt := &RunOpt{}
	for i := range opts {
		opts[i](opt)
	}

	initCtx := opt.InitContext
	if initCtx == nil {
		// default context is not contain http info
		initCtx = NewContext(nil)
	}

	handlerMw := NewHandlerMiddleware(handler)
	if p.globalOrchestrator != nil {
		p.mws = p.globalOrchestrator(p.mws)
	}
	if p.mwOrchestrator != nil {
		p.mws = p.mwOrchestrator(p.mws)
	}
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
	return initCtx.Output(), err
}
