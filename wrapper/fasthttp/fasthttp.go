package fasthttp

import (
	"net/http"

	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/wrapper"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type HttpRespWriterWrapper struct {
	raw *fasthttp.Response
}

func (r *HttpRespWriterWrapper) SetHeader(key, val string) {
	r.raw.Header.Set(key, val)
}

func (r *HttpRespWriterWrapper) GetHeader(key string) string {
	return string(r.raw.Header.Peek(key))
}

func (r *HttpRespWriterWrapper) GetHeaderValues(key string) []string {
	// fasthttp does not support multiple values, plz refer to https://github.com/valyala/fasthttp/issues/179
	return nil
}

func (r *HttpRespWriterWrapper) DelHeader(key string) {
	r.raw.Header.Del(key)
}

func (r *HttpRespWriterWrapper) Write(bs []byte) (int, error) {
	return r.raw.BodyWriter().Write(bs)
}

func (r *HttpRespWriterWrapper) WriteHeader(statusCode int) {
	r.raw.SetStatusCode(statusCode)
}

func (r *HttpRespWriterWrapper) StdHttpWriter() http.ResponseWriter {
	return nil
}

func Wraps(handler core.Handler, opts ...wrapper.SetWrapOpt) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		newReq := http.Request{}
		if err := fasthttpadaptor.ConvertRequest(ctx, &newReq, true); err != nil {
			panic(err)
		}

		pathFunc := func(key string) string {
			v := ctx.UserValue(key)
			str, ok := v.(string)
			if !ok {
				return ""
			}

			return str
		}

		respWrapper := &HttpRespWriterWrapper{&ctx.Response}
		wrapper.HandleHttpInPipeline(wrapper.HandleHttpInPipelineInput{
			Req:            &newReq,
			RespWriter:     respWrapper,
			PathParamsFunc: pathFunc,
			Handler:        handler,
			Opts:           opts,
		})
	}
}
