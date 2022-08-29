package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInput struct {
	Field string
}

type firstMiddleWare struct {
	next Middleware
}

func (f *firstMiddleWare) SetNext(next Middleware) {
	f.next = next
}

func (f *firstMiddleWare) Handle(context Context) error {
	context.SetInput("fi")
	f.next.Handle(context)
	context.SetOutput("fo")
	return nil
}

type secondMiddleWare struct {
	next Middleware
}

func (f *secondMiddleWare) SetNext(next Middleware) {
	f.next = next
}

func (f *secondMiddleWare) Handle(context Context) error {
	context.SetInput("si")
	f.next.Handle(context)
	context.SetOutput("so")
	return nil
}

func TestPipeWork(t *testing.T) {
	input := "hello"
	resp, err := NewPipe(nil).
		AddIf(&firstMiddleWare{}, false).
		Add(&secondMiddleWare{}).
		Run(func(ctx Context) (interface{}, error) {
			assert.Equal(t, "si", ctx.Input())
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, "so", resp)

	resp, err = NewPipe(nil).
		AddIf(&firstMiddleWare{}, true).
		Add(&secondMiddleWare{}).
		Run(func(ctx Context) (interface{}, error) {
			assert.Equal(t, "si", ctx.Input())
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, "fo", resp)

	resp, err = NewPipe(nil).
		AddIf(&firstMiddleWare{}, true).
		AddIf(&secondMiddleWare{}, false).
		Run(func(ctx Context) (interface{}, error) {
			assert.Equal(t, "fi", ctx.Input())
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, "fo", resp)

	// Orchestrator case
	order := 0
	gOrc := func(mws []Middleware) []Middleware {
		assert.Equal(t, 0, order)
		order++
		return mws
	}
	resp, err = NewPipe(gOrc).
		Add(&firstMiddleWare{}).
		SetOrchestrator(func(mws []Middleware) []Middleware {
			assert.Equal(t, 1, len(mws))
			assert.Equal(t, 1, order)
			order++
			return nil
		}).
		Run(func(ctx Context) (interface{}, error) {
			assert.Equal(t, nil, ctx.Input())
			return input, nil
		})

}
