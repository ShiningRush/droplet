package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
)

func Wraps(handler droplet.Handler) func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}