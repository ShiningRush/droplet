package droplet

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	return f.next.Handle(context)
}

type changeResultMiddleWare struct {
	next Middleware
}

func (f *changeResultMiddleWare) SetNext(next Middleware) {
	f.next = next
}

func (f *changeResultMiddleWare) Handle(context Context) error {
	f.next.Handle(context)
	context.SetOutput("not good")
	return nil
}

func TestAddIf(t *testing.T) {
	input := "hello"
	resp, err := NewPipe().
		Add(&firstMiddleWare{}).
		AddIf(&changeResultMiddleWare{}, false).
		Run(func(ctx Context) (interface{}, error) {
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, "hello", resp)

	resp, err = NewPipe().
		Add(&firstMiddleWare{}).
		AddIf(&changeResultMiddleWare{}, true).
		Run(func(ctx Context) (interface{}, error) {
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, "not good", resp)
}

func TestChangeResultPipeWork(t *testing.T) {
	input := "hello"
	resp, err := NewPipe().
		Add(&firstMiddleWare{}).
		Add(&changeResultMiddleWare{}).
		Run(func(ctx Context) (interface{}, error) {
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, "not good", resp)
}

func TestPipeWork(t *testing.T) {
	input := "hello"
	resp, err := NewPipe().
		Add(&firstMiddleWare{}).
		Run(func(ctx Context) (interface{}, error) {
			return input, nil
		})

	assert.NoError(t, err)
	assert.Equal(t, input, resp)
}
