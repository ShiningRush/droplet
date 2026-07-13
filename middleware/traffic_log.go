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
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/log"
)

type TrafficLogMiddleware struct {
	opt *TrafficLogOpt
	BaseMiddleware
}

type RequestTrafficLog struct {
	Context context.Context `json:"-"`

	TraceID   string      `json:"trace_id,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Path      string      `json:"path,omitempty"`
	Method    string      `json:"method,omitempty"`
	Input     interface{} `json:"request,omitempty"`
}

type ResponseTrafficLog struct {
	Context context.Context `json:"-"`

	TraceID     string      `json:"trace_id,omitempty"`
	RequestID   string      `json:"request_id,omitempty"`
	Path        string      `json:"path,omitempty"`
	Method      string      `json:"method,omitempty"`
	ElapsedTime int64       `json:"duration_ms,omitempty"`
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

const missingTrafficLogFieldValue = "-"

type defaultTrafficLogger struct {
}

func trafficLogContextString(ctx context.Context, key string) string {
	value, _ := ctx.Value(key).(string)
	return value
}

func (l *defaultTrafficLogger) LogRequest(tr *RequestTrafficLog) {
	fields := []string{
		"path",
		"method",
	}
	values := []any{
		tr.Path,
		tr.Method,
	}
	appendContextField := func(field string, key string) {
		value := trafficLogContextString(tr.Context, key)
		if value == "" {
			value = missingTrafficLogFieldValue
		}
		fields = append(fields, field)
		values = append(values, value)
	}
	appendContextField("operator", KeyOperator)
	appendContextField("tenant_id", KeyTenantID)
	appendContextField("app_id", KeyAppID)
	appendContextField("locale", KeyLocale)
	appendContextField("remaining_timeout_ms", KeyRemainingTimeoutMS)
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
		"path",
		"method",
		"duration_ms",
	}
	values := []any{
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

	if tr.Error != nil {
		l.logFailedRequest(tr)
	}
	log.CtxInfof(tr.Context, pFmt.String(), values...)
}

func (l *defaultTrafficLogger) logFailedRequest(tr *ResponseTrafficLog) {
	code := data.CodeOf(tr.Error)
	fields := []string{
		"path",
		"method",
		"duration_ms",
		"code",
		"err",
	}
	values := []any{
		tr.Path,
		tr.Method,
		tr.ElapsedTime,
		code,
		tr.Error,
	}

	pFmt := strings.Builder{}
	pFmt.WriteString("request failed, ")
	for _, field := range fields {
		pFmt.WriteString(fmt.Sprintf("%s: %%v, ", field))
	}

	if isWarnErrorCode(code) {
		log.CtxWarnf(tr.Context, pFmt.String(), values...)
		return
	}
	log.CtxErrorf(tr.Context, pFmt.String(), values...)
}

func isWarnErrorCode(code int) bool {
	return code != 0 && code != data.ErrCodeInternal
}

func (mw *TrafficLogMiddleware) Handle(ctx core.Context) error {
	if mw.opt.Logger == nil {
		mw.opt.Logger = defaultLogger
	}

	reqLog := &RequestTrafficLog{
		Context:   ctx.Context(),
		TraceID:   ctx.GetString(KeyTraceID),
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
		TraceID:   ctx.GetString(KeyTraceID),
		RequestID: ctx.GetString(KeyRequestID),
		Path:      ctx.Path(),
		Method:    ctx.Get(KeyHttpRequest).(*http.Request).Method,
	}
	now := time.Now()
	respLog.Error = mw.BaseMiddleware.Handle(ctx)
	respLog.ElapsedTime = time.Since(now).Milliseconds()
	setAppServerTiming(ctx.ResponseHeader(), time.Duration(respLog.ElapsedTime)*time.Millisecond)

	if mw.opt.LogResp {
		respLog.Output = ctx.Output()
	}

	mw.opt.Logger.LogResponse(respLog)

	return respLog.Error
}
