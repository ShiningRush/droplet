//go:generate mockgen -source traffic_log.go  -destination traffic_log_mock.go -package middleware
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Path        string      `json:"path,omitempty"`
	Method      string      `json:"method,omitempty"`
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
	fields := []string{
		"request_id",
		"path",
		"method",
	}
	values := []any{
		tr.RequestID,
		tr.Path,
		tr.Method,
	}
	if tr.Input != nil {
		fields = append(fields, "input")
		input, _ := json.Marshal(tr.Input)
		values = append(values, string(input))
	}

	pFmt := strings.Builder{}
	pFmt.WriteString("request start, ")
	for _, field := range fields {
		pFmt.WriteString(fmt.Sprintf("%s: %%v, ", field))
	}

	log.CtxInfof(tr.Context, pFmt.String(), values...)
}

func (l *defaultTrafficLogger) LogResponse(tr *ResponseTrafficLog) {
	fields := []string{
		"request_id",
		"path",
		"method",
		"elapsed_time",
	}
	values := []any{
		tr.RequestID,
		tr.Path,
		tr.Method,
		tr.ElapsedTime,
	}

	if tr.Error != nil {
		fields = append(fields, "err")
		values = append(values, tr.Error)
	}
	if tr.Output != nil {
		fields = append(fields, "output")
		output, _ := json.Marshal(tr.Output)
		values = append(values, output)
	}

	pFmt := strings.Builder{}
	pFmt.WriteString("request complete, ")
	for _, field := range fields {
		pFmt.WriteString(fmt.Sprintf("%s: %%v, ", field))
	}

	log.CtxInfof(tr.Context, pFmt.String(), values...)
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
		Path:      ctx.Path(),
		Method:    ctx.Get(KeyHttpRequest).(*http.Request).Method,
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
