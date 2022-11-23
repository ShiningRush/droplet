package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shiningrush/droplet/core"
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

			defaultLogger = ml
			err := testMw.Handle(c)
			assert.Equal(t, tc.wantRespLog.Error, err)
		})
	}
}
