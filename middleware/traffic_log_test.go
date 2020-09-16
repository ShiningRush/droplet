package middleware

import (
	"fmt"
	"github.com/shiningrush/droplet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

func TestTrafficLogMiddleware_Handle(t *testing.T) {
	tests := []struct {
		giveOpt TrafficLogOpt
		wantLog *TrafficLog
	}{
		{
			giveOpt: TrafficLogOpt{
				IsLogReqAndResp: true,
			},
			wantLog: &TrafficLog{
				Path:        "path",
				Method:      "method",
				ElapsedTime: 100,
				Input:       "req",
				Output:      "resp",
				Error:       fmt.Errorf("failed"),
			},
		},
	}

	for _, tc := range tests {
		mMw := &droplet.MockMiddleware{}
		mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
			time.Sleep(time.Duration(tc.wantLog.ElapsedTime) * time.Millisecond)
		}).Return(tc.wantLog.Error)

		testMw := TrafficLogMiddleware{
			opt: tc.giveOpt,
			BaseMiddleware: BaseMiddleware{
				next: mMw,
			},
		}
		c := droplet.NewContext()
		c.SetPath(tc.wantLog.Path)
		c.Set(KeyHttpRequest, &http.Request{
			Method: tc.wantLog.Method,
		})
		c.SetInput(tc.wantLog.Input)
		c.SetOutput(tc.wantLog.Output)
		testMw.opt.LogFunc = func(log *TrafficLog) {
			assert.Equal(t, tc.wantLog, log)
		}
		err := testMw.Handle(c)
		assert.Equal(t, tc.wantLog.Error, err)
	}
}
