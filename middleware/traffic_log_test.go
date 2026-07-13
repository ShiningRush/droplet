package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTrafficLogMiddleware_Handle(t *testing.T) {
	tests := []struct {
		name        string
		giveOpt     TrafficLogOpt
		wantReqLog  *RequestTrafficLog
		wantRespLog *ResponseTrafficLog
	}{
		{
			name: "normal",
			giveOpt: TrafficLogOpt{
				LogReq:  true,
				LogResp: true,
			},
			wantReqLog: &RequestTrafficLog{
				TraceID:   "trace",
				RequestID: "req",
				Path:      "path",
				Method:    "method",
				Input:     "req",
			},
			wantRespLog: &ResponseTrafficLog{
				TraceID:   "trace",
				RequestID: "req",
				Path:      "path",
				Method:    "method",
				Error:     fmt.Errorf("test errr"),
				Output:    "output",
			},
		},
		{
			name: "do not log req resp",
			giveOpt: TrafficLogOpt{
				LogReq:  false,
				LogResp: false,
			},
			wantReqLog: &RequestTrafficLog{
				TraceID:   "trace",
				RequestID: "req",
				Path:      "path",
				Method:    "method",
			},
			wantRespLog: &ResponseTrafficLog{
				TraceID:   "trace",
				RequestID: "req",
				Path:      "path",
				Method:    "method",
				Error:     fmt.Errorf("test errr"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mMw := &core.MockMiddleware{}
			mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
				time.Sleep(time.Duration(tc.wantRespLog.ElapsedTime) * time.Millisecond)
			}).Return(tc.wantRespLog.Error)

			testMw := TrafficLogMiddleware{
				opt: &tc.giveOpt,
				BaseMiddleware: BaseMiddleware{
					next: mMw,
				},
			}
			c := core.NewContext()
			c.SetPath(tc.wantReqLog.Path)
			c.Set(KeyHttpRequest, &http.Request{
				Method: tc.wantReqLog.Method,
			})
			c.Set(KeyTraceID, "trace")
			c.Set(KeyRequestID, "req")
			c.SetContext(contextWithRequestID(t, "trace", "req"))
			c.SetInput(tc.wantReqLog.Input)
			c.SetOutput(tc.wantRespLog.Output)
			tc.wantReqLog.Context = c.Context()
			tc.wantRespLog.Context = c.Context()

			ml := NewMockTrafficLogger(ctrl)
			ml.EXPECT().LogRequest(gomock.Any()).Do(func(reqLog interface{}) {
				assert.Equal(t, tc.wantReqLog, reqLog)
			})
			ml.EXPECT().LogResponse(gomock.Any()).Do(func(respLog interface{}) {
				assert.Equal(t, tc.wantRespLog, respLog)
			})

			oldDefaultLogger := defaultLogger
			defaultLogger = ml
			defer func() {
				defaultLogger = oldDefaultLogger
			}()
			err := testMw.Handle(c)
			assert.Equal(t, tc.wantRespLog.Error, err)
		})
	}
}

func TestDefaultTrafficLogger_LogRequestUsesContextFields(t *testing.T) {
	logger := &captureLogger{}
	oldLogger := log.DefLogger
	log.DefLogger = logger
	defer func() {
		log.DefLogger = oldLogger
	}()

	(&defaultTrafficLogger{}).LogRequest(&RequestTrafficLog{
		Context:   contextWithRequestFields(t),
		TraceID:   "trace",
		RequestID: "req",
		Path:      "/path",
		Method:    http.MethodGet,
	})

	assert.Len(t, logger.entries, 1)
	assert.Equal(t, "info", logger.entries[0].level)
	assert.True(t, strings.HasPrefix(logger.entries[0].msg, "request start, "))
	assert.NotContains(t, logger.entries[0].msg, "trace_id")
	assert.NotContains(t, logger.entries[0].msg, "request_id")
	assert.NotContains(t, logger.entries[0].args, "trace")
	assert.NotContains(t, logger.entries[0].args, "req")
	assert.Contains(t, logger.entries[0].msg, "operator")
	assert.Contains(t, logger.entries[0].msg, "tenant_id")
	assert.Contains(t, logger.entries[0].msg, "app_id")
	assert.Contains(t, logger.entries[0].msg, "locale")
	assert.Contains(t, logger.entries[0].msg, "remaining_timeout_ms")
	assert.Contains(t, logger.entries[0].args, "operator")
	assert.Contains(t, logger.entries[0].args, "tenant")
	assert.Contains(t, logger.entries[0].args, "app")
	assert.Contains(t, logger.entries[0].args, "zh-CN")
	assert.Contains(t, logger.entries[0].args, "1000")
	assert.Equal(t, "trace", logger.entries[0].ctx.Value(KeyTraceID))
	assert.Equal(t, "req", logger.entries[0].ctx.Value(KeyRequestID))
}

func TestDefaultTrafficLogger_LogRequestUsesPlaceholderForMissingContextFields(t *testing.T) {
	logger := &captureLogger{}
	oldLogger := log.DefLogger
	log.DefLogger = logger
	defer func() {
		log.DefLogger = oldLogger
	}()

	(&defaultTrafficLogger{}).LogRequest(&RequestTrafficLog{
		Context: contextWithRequestID(t, "trace", "req"),
		Path:    "/path",
		Method:  http.MethodGet,
	})

	assert.Len(t, logger.entries, 1)
	assert.Contains(t, logger.entries[0].msg, "operator")
	assert.Contains(t, logger.entries[0].msg, "tenant_id")
	assert.Contains(t, logger.entries[0].msg, "app_id")
	assert.Contains(t, logger.entries[0].msg, "locale")
	assert.Contains(t, logger.entries[0].msg, "remaining_timeout_ms")

	placeholderCount := 0
	for _, arg := range logger.entries[0].args {
		if arg == missingTrafficLogFieldValue {
			placeholderCount++
		}
	}
	assert.Equal(t, 5, placeholderCount)
}

func TestDefaultTrafficLogger_LogResponseLogsFailedRequestLevel(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantLevel string
		wantCode  int
	}{
		{
			name:      "validate error logs warn",
			err:       data.NewValidateError("invalid request", nil),
			wantLevel: "warn",
			wantCode:  data.ErrCodeValidate,
		},
		{
			name:      "friendly error logs warn",
			err:       data.NewFriendlyError("friendly error"),
			wantLevel: "warn",
			wantCode:  data.ErrCodeFriendly,
		},
		{
			name: "custom coded error logs warn",
			err: &codedTestError{
				code:    20001,
				message: "business error",
			},
			wantLevel: "warn",
			wantCode:  20001,
		},
		{
			name:      "internal error logs error",
			err:       data.NewInternalError("internal error"),
			wantLevel: "error",
			wantCode:  data.ErrCodeInternal,
		},
		{
			name:      "plain error logs error",
			err:       errors.New("plain error"),
			wantLevel: "error",
			wantCode:  data.ErrCodeInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := &captureLogger{}
			oldLogger := log.DefLogger
			log.DefLogger = logger
			defer func() {
				log.DefLogger = oldLogger
			}()

			(&defaultTrafficLogger{}).LogResponse(&ResponseTrafficLog{
				Context:     contextWithRequestID(t, "trace", "req"),
				TraceID:     "trace",
				RequestID:   "req",
				Path:        "/path",
				Method:      http.MethodGet,
				ElapsedTime: 10,
				Error:       tc.err,
			})

			assert.Len(t, logger.entries, 2)
			assert.Equal(t, tc.wantLevel, logger.entries[0].level)
			assert.True(t, strings.HasPrefix(logger.entries[0].msg, "request failed, "))
			assert.NotContains(t, logger.entries[0].msg, "trace_id")
			assert.NotContains(t, logger.entries[0].msg, "request_id")
			assert.Contains(t, logger.entries[0].msg, "duration_ms")
			assert.NotContains(t, logger.entries[0].args, "trace")
			assert.NotContains(t, logger.entries[0].args, "req")
			assert.Equal(t, "trace", logger.entries[0].ctx.Value(KeyTraceID))
			assert.Equal(t, "req", logger.entries[0].ctx.Value(KeyRequestID))
			assert.Contains(t, logger.entries[0].args, tc.wantCode)
			assert.Contains(t, logger.entries[0].args, tc.err)
			assert.Equal(t, "info", logger.entries[1].level)
			assert.True(t, strings.HasPrefix(logger.entries[1].msg, "request complete, "))
			assert.NotContains(t, logger.entries[1].msg, "trace_id")
			assert.NotContains(t, logger.entries[1].msg, "request_id")
			assert.Contains(t, logger.entries[1].msg, "duration_ms")
			assert.NotContains(t, logger.entries[1].args, "trace")
			assert.NotContains(t, logger.entries[1].args, "req")
			assert.Equal(t, "trace", logger.entries[1].ctx.Value(KeyTraceID))
			assert.Equal(t, "req", logger.entries[1].ctx.Value(KeyRequestID))
		})
	}
}

func TestResponseTrafficLogJSONUsesDurationMS(t *testing.T) {
	bs, err := json.Marshal(&ResponseTrafficLog{
		ElapsedTime: 10,
	})
	assert.NoError(t, err)
	assert.Contains(t, string(bs), `"duration_ms":10`)
	assert.NotContains(t, string(bs), `"elapsed_time"`)
}

func TestTrafficLogMiddlewareUsesElapsedTimeForServerTiming(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mMw := &core.MockMiddleware{}
	mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
		time.Sleep(10 * time.Millisecond)
	}).Return(nil)

	var capturedRespLog *ResponseTrafficLog
	ml := NewMockTrafficLogger(ctrl)
	ml.EXPECT().LogRequest(gomock.Any())
	ml.EXPECT().LogResponse(gomock.Any()).Do(func(respLog interface{}) {
		capturedRespLog = respLog.(*ResponseTrafficLog)
	})

	testMw := TrafficLogMiddleware{
		opt: &TrafficLogOpt{
			Logger: ml,
		},
		BaseMiddleware: BaseMiddleware{
			next: mMw,
		},
	}

	c := core.NewContext()
	c.SetPath("path")
	c.Set(KeyHttpRequest, &http.Request{
		Method: http.MethodGet,
	})
	c.Set(KeyTraceID, "trace")
	c.Set(KeyRequestID, "req")
	c.SetContext(contextWithRequestID(t, "trace", "req"))

	err := testMw.Handle(c)
	assert.NoError(t, err)
	if assert.NotNil(t, capturedRespLog) {
		assert.Equal(t, formatServerTiming(time.Duration(capturedRespLog.ElapsedTime)*time.Millisecond), c.ResponseHeader().Get(serverTimingHeader))
	}
}

func contextWithRequestID(t *testing.T, traceID string, requestID string) context.Context {
	t.Helper()

	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyTraceID, traceID)
	ctx = context.WithValue(ctx, KeyRequestID, requestID)
	return ctx
}

func contextWithRequestFields(t *testing.T) context.Context {
	t.Helper()

	ctx := contextWithRequestID(t, "trace", "req")
	ctx = context.WithValue(ctx, KeyOperator, "operator")
	ctx = context.WithValue(ctx, KeyTenantID, "tenant")
	ctx = context.WithValue(ctx, KeyAppID, "app")
	ctx = context.WithValue(ctx, KeyLocale, "zh-CN")
	ctx = context.WithValue(ctx, KeyRemainingTimeoutMS, "1000")
	return ctx
}

type logEntry struct {
	level string
	ctx   context.Context
	msg   string
	args  []interface{}
}

type captureLogger struct {
	entries []logEntry
}

func (l *captureLogger) append(level string, ctx context.Context, msg string, args ...interface{}) {
	l.entries = append(l.entries, logEntry{
		level: level,
		ctx:   ctx,
		msg:   msg,
		args:  args,
	})
}

func (l *captureLogger) CtxDebugf(ctx context.Context, msg string, args ...interface{}) {
	l.append("debug", ctx, msg, args...)
}

func (l *captureLogger) CtxInfof(ctx context.Context, msg string, args ...interface{}) {
	l.append("info", ctx, msg, args...)
}

func (l *captureLogger) CtxWarnf(ctx context.Context, msg string, args ...interface{}) {
	l.append("warn", ctx, msg, args...)
}

func (l *captureLogger) CtxErrorf(ctx context.Context, msg string, args ...interface{}) {
	l.append("error", ctx, msg, args...)
}

func (l *captureLogger) CtxFatalf(ctx context.Context, msg string, args ...interface{}) {
	l.append("fatal", ctx, msg, args...)
}

func (l *captureLogger) Debugf(msg string, args ...interface{}) {
	l.append("debug", nil, msg, args...)
}

func (l *captureLogger) Infof(msg string, args ...interface{}) {
	l.append("info", nil, msg, args...)
}

func (l *captureLogger) Warnf(msg string, args ...interface{}) {
	l.append("warn", nil, msg, args...)
}

func (l *captureLogger) Errorf(msg string, args ...interface{}) {
	l.append("error", nil, msg, args...)
}

func (l *captureLogger) Fatalf(msg string, args ...interface{}) {
	l.append("fatal", nil, msg, args...)
}
