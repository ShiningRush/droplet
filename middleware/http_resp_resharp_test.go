package middleware

import (
	"fmt"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestHttpRespReshapeMiddleware_Handle(t *testing.T) {
	tests := []struct {
		giveResp interface{}
		giveErr  error
		wantErr  error
		wantResp interface{}
	}{
		{
			giveErr: fmt.Errorf("failed"),
			wantResp: &data.Response{
				Code:    data.ErrCodeInternal,
				Message: "failed",
			},
		},
		{
			giveErr: fmt.Errorf("failed"),
			giveResp: &data.Response{
				Code:    http.StatusOK,
				Message: "OK",
			},
			wantResp: &data.Response{
				Code:    data.ErrCodeInternal,
				Message: "failed",
			},
		},
		{
			giveErr: fmt.Errorf("failed"),
			giveResp: &data.SpecCodeResponse{
				Response: data.Response{
					Code:    http.StatusOK,
					Message: "OK",
				},
			},
			wantResp: &data.SpecCodeResponse{
				Response: data.Response{
					Code:    data.ErrCodeInternal,
					Message: "failed",
				},
			},
		},
		{
			giveErr: &data.BaseError{
				Code:    data.ErrCodeFriendly,
				Message: "friendly error",
			},
			wantResp: &data.Response{
				Code:    data.ErrCodeFriendly,
				Message: "friendly error",
			},
		},
		{
			giveResp: "test",
			wantResp: &data.Response{
				Data: "test",
			},
		},
		{
			giveResp: &data.Response{
				Data: "test",
			},
			wantResp: &data.Response{
				Data: "test",
			},
		},
	}

	for _, tc := range tests {
		mMw := &droplet.MockMiddleware{}
		mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
			ctx := args.Get(0).(droplet.Context)
			ctx.SetOutput(tc.giveResp)
		}).Return(tc.giveErr)

		testMw := HttpRespReshapeMiddleware{
			BaseMiddleware{
				next: mMw,
			},
		}
		c := droplet.NewContext()
		err := testMw.Handle(c)
		if err != nil {
			assert.Equal(t, tc.wantErr, err)
			continue
		}
		assert.Equal(t, tc.wantResp, c.Output())
	}
}
