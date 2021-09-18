package gorestful

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/log"
	"github.com/shiningrush/droplet/middleware"
	"github.com/shiningrush/droplet/wrapper"
	"net/http"
)

func Wraps(handler droplet.Handler, opts ...wrapper.SetWrapOpt) func(*restful.Request,*restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		opt := &wrapper.WrapOptBase{}
		for i := range opts {
			opts[i](opt)
		}

		dCtx := droplet.NewContext()
		dCtx.SetContext(request.Request.Context())

		ret, _ := droplet.NewPipe().
			Add(middleware.NewHttpInfoInjectorMiddleware(middleware.HttpInfoInjectorOption{
				ReqFunc: func() *http.Request {
					return request.Request
				},
			})).
			Add(middleware.NewRespReshapeMiddleware()).
			Add(middleware.NewHttpInputMiddleWare(middleware.HttpInputOption{
				PathParamsFunc: func(key string) string {
					return request.PathParameter(key)
				},
				InputType:      opt.InputType,
				IsReadFromBody: opt.IsReadFromBody,
			})).
			Add(middleware.NewTrafficLogMiddleware(opt.TrafficLogOpt)).
			SetOrchestrator(opt.Orchestrator).
			Run(handler, droplet.InitContext(dCtx))

		switch ret.(type) {
		case *data.FileResponse:
			fr := ret.(*data.FileResponse)
			if fr.ContentType == "" {
				fr.ContentType = "application/octet-stream"
			}
			response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fr.Name))
			response.Header().Set("Content-type", fr.ContentType)
			_, err := response.Write(fr.Content)
			if err != nil {
				log.Error("write response failed",
					"err", err)
			}
		case droplet.SpecCodeHttpResponse:
			resp := ret.(droplet.SpecCodeHttpResponse)
			if err := response.WriteHeaderAndJson(resp.GetStatusCode(), resp, "application/json"); err != nil {
				log.Error("write resp failed",
					"err", err,
					"path", request.Request.URL.Path)
			}
		default:
			if err := response.WriteHeaderAndJson(http.StatusOK, ret, "application/json"); err != nil {
				log.Error("write resp failed",
					"err", err,
					"path", request.Request.URL.Path)
			}
		}
	}
}
