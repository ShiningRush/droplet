package middleware

import (
	"encoding/json"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/log"
	"net/http"
	"time"
)

type TrafficLogMiddleware struct {
	opt TrafficLogOpt
	BaseMiddleware
}

type TrafficLog struct {
	RequestID   string      `json:"request_id,omitempty"`
	Path        string      `json:"path,omitempty"`
	Method      string      `json:"method,omitempty"`
	ElapsedTime int64       `json:"elapsed_time,omitempty"`
	Input       interface{} `json:"request,omitempty"`
	Output      interface{} `json:"response,omitempty"`
	Error       error       `json:"error,omitempty"`
}

type TrafficLogOpt struct {
	IsLogReqAndResp bool
	LogFunc         func(log *TrafficLog)
}

func NewTrafficLogMiddleware(opt TrafficLogOpt) *TrafficLogMiddleware {
	return &TrafficLogMiddleware{
		opt: opt,
	}
}

var defaultLogFunc = func(tl *TrafficLog) {
	if tl.Input == nil && tl.Output == nil {
		log.Info("request finished",
			"request_id", tl.RequestID,
			"path", tl.Path,
			"method", tl.Method,
			"elapsed_time", tl.ElapsedTime,
			"err", tl.Error)
		return
	}
	input, _ := json.Marshal(tl.Input)
	output, _ := json.Marshal(tl.Output)
	log.Info("request finished",
		"request_id", tl.RequestID,
		"path", tl.Path,
		"method", tl.Method,
		"elapsed_time", tl.ElapsedTime,
		"err", tl.Error,
		"input", input,
		"output", output)
}

func (mw *TrafficLogMiddleware) Handle(ctx droplet.Context) error {
	if mw.opt.LogFunc == nil {
		mw.opt.LogFunc = defaultLogFunc
	}

	logMsg := &TrafficLog{
		RequestID: ctx.GetString(KeyRequestID),
		Path:      ctx.Path(),
		Method:    ctx.Get(KeyHttpRequest).(*http.Request).Method,
	}
	logMsg.Path = ctx.Path()
	logMsg.Method = ctx.Get(KeyHttpRequest).(*http.Request).Method

	now := time.Now()
	logMsg.Error = mw.BaseMiddleware.Handle(ctx)
	logMsg.ElapsedTime = time.Since(now).Nanoseconds() / 1000 / 1000 // ns to ms

	if mw.opt.IsLogReqAndResp {
		logMsg.Input = ctx.Input()
		logMsg.Output = ctx.Output()
	}

	mw.opt.LogFunc(logMsg)

	return logMsg.Error
}
