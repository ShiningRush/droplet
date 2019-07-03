package middlewares

import (
	"github.com/shiningrush/droplet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInputMiddleWare(t *testing.T) {
	input := "hello"
	_, err := droplet.NewPipe().
		Add(NewInputMiddleWare(input)).
		Run(func(ctx droplet.Context) (interface{}, error) {
			assert.Equal(t, input, ctx.Input(), "ctx.input should equal input")
			return input, nil
		})

	assert.NoError(t, err, "pipe should no error")
}
