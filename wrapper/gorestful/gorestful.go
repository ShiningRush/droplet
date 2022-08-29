package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/wrapper"
)

func Wraps(handler core.Handler, opts ...wrapper.SetWrapOpt) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		wrapper.HandleHttpInPipeline(wrapper.HandleHttpInPipelineInput{
			Req:            request.Request,
			RespWriter:     wrapper.NewResponseWriter(response.ResponseWriter),
			PathParamsFunc: request.PathParameter,
			Handler:        handler,
			Opts:           opts,
		})
	}
}
