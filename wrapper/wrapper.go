package wrapper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/log"
	"github.com/shiningrush/droplet/middleware"
)

type HandleHttpInPipelineInput struct {
	Req            *http.Request
	RespWriter     droplet.ResponseWriter
	PathParamsFunc func(key string) string
	Handler        droplet.Handler
	Opts           []SetWrapOpt
}

func HandleHttpInPipeline(input HandleHttpInPipelineInput) {
	opt := &WrapOptBase{}
	for _, op := range input.Opts {
		op(opt)
	}

	dCtx := droplet.NewContext()
	dCtx.SetContext(input.Req.Context())

	ret, _ := droplet.NewPipe().
		Add(middleware.NewHttpInfoInjectorMiddleware(middleware.HttpInfoInjectorOption{
			ReqFunc: func() *http.Request {
				return input.Req
			},
		})).
		Add(middleware.NewRespReshapeMiddleware()).
		Add(middleware.NewHttpInputMiddleWare(middleware.HttpInputOption{
			PathParamsFunc: input.PathParamsFunc,
			InputType:      opt.InputType,
			IsReadFromBody: opt.IsReadFromBody,
		})).
		Add(middleware.NewTrafficLogMiddleware(opt.TrafficLogOpt)).
		SetOrchestrator(opt.Orchestrator).
		Run(input.Handler, droplet.InitContext(dCtx))

	switch ret.(type) {
	case droplet.RawHttpResponse:
		rr := ret.(droplet.RawHttpResponse)
		if err := rr.WriteRawResponse(input.RespWriter); err != nil {
			logWriteErrors(input.Req, err)
		}
	case droplet.HttpFileResponse:
		fr := ret.(droplet.HttpFileResponse)
		name, contentType, content := fr.Get()
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		input.RespWriter.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
		input.RespWriter.SetHeader("Content-type", contentType)
		_, err := input.RespWriter.Write(content)
		if err != nil {
			logWriteErrors(input.Req, err)
		}
	case droplet.SpecCodeHttpResponse:
		resp := ret.(droplet.SpecCodeHttpResponse)
		if err := writeJsonToResp(input.RespWriter, resp.GetStatusCode(), resp); err != nil {
			logWriteErrors(input.Req, err)
		}
	default:
		if err := writeJsonToResp(input.RespWriter, http.StatusOK, ret); err != nil {
			logWriteErrors(input.Req, err)
		}
	}
}

func logWriteErrors(req *http.Request, err error) {
	log.Error("write resp failed",
		"err", err,
		"url", req.URL.String())
}

func writeJsonToResp(rw droplet.ResponseWriter, code int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	rw.SetHeader("Content-Type", "application/json")

	rw.WriteHeader(code)
	if _, err := rw.Write(bs); err != nil {
		return fmt.Errorf("write to response failed: %w", err)
	}
	return nil
}
