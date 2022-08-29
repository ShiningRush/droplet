package main

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/wrapper"
	ginwrap "github.com/shiningrush/droplet/wrapper/gin"
)

func main() {
	r := gin.Default()
	r.POST("/json_input/:id", ginwrap.Wraps(JsonInputDo, wrapper.InputType(reflect.TypeOf(&JsonInput{}))))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type JsonInput struct {
	ID    string   `auto_read:"id,path" json:"id"`
	User  string   `auto_read:"user,header" json:"user"`
	IPs   []string `json:"ips"`
	Count int      `json:"count"`
	Body  []byte   `auto_read:"@body"`
}

func JsonInputDo(ctx core.Context) (interface{}, error) {
	input := ctx.Input().(*JsonInput)

	return input, nil
}
