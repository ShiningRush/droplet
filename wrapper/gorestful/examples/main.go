package main

import (
	"net/http"
	"reflect"

	"github.com/emicklei/go-restful/v3"
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/wrapper"

	rwrap "github.com/shiningrush/droplet/wrapper/gorestful"
)

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.POST("/json_input/{id}").To(rwrap.Wraps(JsonInputDo, wrapper.InputType(reflect.TypeOf(&JsonInput{})))))
	restful.Add(ws)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
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
