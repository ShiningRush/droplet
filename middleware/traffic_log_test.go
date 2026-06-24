package middleware

import (
	"context"
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
				RequestID: "req",
				Path:      "path",
				Method:    "method",
				Input:     "req",
			},
			wantRespLog: &ResponseTrafficLog{
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
				RequestID: "req",
				Path:      "path",
				Method:    "method",
			},
			wantRespLog: &ResponseTrafficLog{
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
			c.Set(KeyRequestID, tc.wantReqLog.RequestID)
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
				Context:     contextWithRequestID(t, "req"),
				RequestID:   "req",
				Path:        "/path",
				Method:      http.MethodGet,
				ElapsedTime: 10,
				Error:       tc.err,
			})

			assert.Len(t, logger.entries, 2)
			assert.Equal(t, tc.wantLevel, logger.entries[0].level)
			assert.True(t, strings.HasPrefix(logger.entries[0].msg, "request failed, "))
			assert.Contains(t, logger.entries[0].args, tc.wantCode)
			assert.Contains(t, logger.entries[0].args, tc.err)
			assert.Equal(t, "info", logger.entries[1].level)
			assert.True(t, strings.HasPrefix(logger.entries[1].msg, "request complete, "))
		})
	}
}

func contextWithRequestID(t *testing.T, requestID string) context.Context {
	t.Helper()

	c := core.NewContext()
	c.Set(KeyRequestID, requestID)
	return c.Context()
}

type logEntry struct {
	level string
	msg   string
	args  []interface{}
}

type captureLogger struct {
	entries []logEntry
}

func (l *captureLogger) append(level string, msg string, args ...interface{}) {
	l.entries = append(l.entries, logEntry{
		level: level,
		msg:   msg,
		args:  args,
	})
}

func (l *captureLogger) CtxDebugf(ctx context.Context, msg string, args ...interface{}) {
	l.append("debug", msg, args...)
}

func (l *captureLogger) CtxInfof(ctx context.Context, msg string, args ...interface{}) {
	l.append("info", msg, args...)
}

func (l *captureLogger) CtxWarnf(ctx context.Context, msg string, args ...interface{}) {
	l.append("warn", msg, args...)
}

func (l *captureLogger) CtxErrorf(ctx context.Context, msg string, args ...interface{}) {
	l.append("error", msg, args...)
}

func (l *captureLogger) CtxFatalf(ctx context.Context, msg string, args ...interface{}) {
	l.append("fatal", msg, args...)
}

func (l *captureLogger) Debugf(msg string, args ...interface{}) {
	l.append("debug", msg, args...)
}

func (l *captureLogger) Infof(msg string, args ...interface{}) {
	l.append("info", msg, args...)
}

func (l *captureLogger) Warnf(msg string, args ...interface{}) {
	l.append("warn", msg, args...)
}

func (l *captureLogger) Errorf(msg string, args ...interface{}) {
	l.append("error", msg, args...)
}

func (l *captureLogger) Fatalf(msg string, args ...interface{}) {
	l.append("fatal", msg, args...)
}
