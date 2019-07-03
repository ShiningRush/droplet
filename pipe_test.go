package droplet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testInput struct {
	Field string
}

type firstMiddleWare struct {
	next MiddleWare
}

func (f *firstMiddleWare) SetNext(next MiddleWare) {
	f.next = next
}

func (f *firstMiddleWare) Handle(context Context) error {
	return f.next.Handle(context)
}

type changeResultMiddleWare struct {
	next MiddleWare
}

func (f *changeResultMiddleWare) SetNext(next MiddleWare) {
	f.next = next
}

func (f *changeResultMiddleWare) Handle(context Context) error {
	f.next.Handle(context)
	context.SetOutput("not good")
	return nil
}

func TestChangeResultPipeWork(t *testing.T) {
	input := "hello"
	resp, err := NewPipe().
		Add(&firstMiddleWare{}).
		Add(&changeResultMiddleWare{}).
		Run(func(ctx Context) (interface{}, error) {
			return input, nil
		})

	assert.NoError(t, err, "pipe should no error")
	assert.Equal(t, "not good", resp, "resp should equal input")
}

func TestPipeWork(t *testing.T) {
	input := "hello"
	resp, err := NewPipe().
		Add(&firstMiddleWare{}).
		Run(func(ctx Context) (interface{}, error) {
			return input, nil
		})

	assert.NoError(t, err, "pipe should no error")
	assert.Equal(t, input, resp, "resp should equal input")
}
