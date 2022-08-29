//go:generate mockgen -source traffic_log.go  -destination traffic_log_mock.go -package middleware
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/log"
)

type TrafficLogMiddleware struct {
	opt *TrafficLogOpt
	BaseMiddleware
}

type RequestTrafficLog struct {
	Context context.Context `json:"-"`

	RequestID string      `json:"request_id,omitempty"`
	Path      string      `json:"path,omitempty"`
	Method    string      `json:"method,omitempty"`
	Input     interface{} `json:"request,omitempty"`
}

type ResponseTrafficLog struct {
	Context context.Context `json:"-"`

	RequestID   string      `json:"request_id,omitempty"`
	ElapsedTime int64       `json:"elapsed_time,omitempty"`
	Output      interface{} `json:"response,omitempty"`
	Error       error       `json:"error,omitempty"`
}

type TrafficLogOpt struct {
	LogReq  bool
	LogResp bool
	Logger  TrafficLogger
}

func NewTrafficLogMiddleware(opt *TrafficLogOpt) *TrafficLogMiddleware {
	return &TrafficLogMiddleware{
		opt: opt,
	}
}

type TrafficLogger interface {
	LogRequest(tr *RequestTrafficLog)
	LogResponse(tr *ResponseTrafficLog)
}

var defaultLogger TrafficLogger = &defaultTrafficLogger{}

type defaultTrafficLogger struct {
}

func (l *defaultTrafficLogger) LogRequest(tr *RequestTrafficLog) {
	fields := []interface{}{
		"request_id", tr.RequestID,
		"path", tr.Path,
		"method", tr.Method,
	}
	if tr.Input != nil {
		input, _ := json.Marshal(tr.Input)
		fields = append(fields, []interface{}{
			"input", string(input),
		}...)
	}
	log.CtxInfo(tr.Context, "request start", fields...)
}

func (l *defaultTrafficLogger) LogResponse(tr *ResponseTrafficLog) {
	fields := []interface{}{
		"request_id", tr.RequestID,
		"elapsed_time", tr.ElapsedTime,
	}
	if tr.Error != nil {
		fields = append(fields, []interface{}{
			"err", tr.Error,
		}...)
	}
	if tr.Output != nil {
		output, _ := json.Marshal(tr.Output)
		fields = append(fields, []interface{}{
			"output", string(output),
		}...)
	}
	log.CtxInfo(tr.Context, "request complete", fields...)
}

func (mw *TrafficLogMiddleware) Handle(ctx core.Context) error {
	if mw.opt.Logger == nil {
		mw.opt.Logger = defaultLogger
	}

	reqLog := &RequestTrafficLog{
		Context:   ctx.Context(),
		RequestID: ctx.GetString(KeyRequestID),
		Path:      ctx.Path(),
		Method:    ctx.Get(KeyHttpRequest).(*http.Request).Method,
	}
	if mw.opt.LogReq {
		reqLog.Input = ctx.Input()
	}
	mw.opt.Logger.LogRequest(reqLog)

	respLog := &ResponseTrafficLog{
		Context:   ctx.Context(),
		RequestID: ctx.GetString(KeyRequestID),
	}
	now := time.Now()
	respLog.Error = mw.BaseMiddleware.Handle(ctx)
	respLog.ElapsedTime = time.Since(now).Nanoseconds() / 1000 / 1000 // ns to ms

	if mw.opt.LogResp {
		respLog.Output = ctx.Output()
	}

	mw.opt.Logger.LogResponse(respLog)

	return respLog.Error
}
