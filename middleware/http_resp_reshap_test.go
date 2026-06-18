package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type codedTestError struct {
	code    int
	message string
	cause   error
}

func (e *codedTestError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *codedTestError) Unwrap() error {
	return e.cause
}

func (e *codedTestError) ErrorCode() int {
	return e.code
}

type extraDataTestError struct {
	Cause error
	Data  interface{}
}

func (e *extraDataTestError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "error with extra data"
}

func (e *extraDataTestError) Unwrap() error {
	return e.Cause
}

func TestHttpRespReshapeMiddleware_Handle(t *testing.T) {
	tests := []struct {
		name        string
		giveOptCode int
		giveResp    interface{}
		giveErr     error
		wantErr     error
		wantResp    interface{}
	}{
		{
			name:    "def-code",
			giveErr: fmt.Errorf("failed"),
			wantResp: &data.Response{
				Code:    data.ErrCodeInternal,
				Message: "failed",
			},
		},
		{
			name:        "opt-code",
			giveOptCode: 500,
			giveErr:     fmt.Errorf("failed"),
			wantResp: &data.Response{
				Code:    500,
				Message: "failed",
			},
		},
		{
			name:    "err-diff",
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
			name:    "spec-status-code",
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
			name: "friend",
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
			name: "wrap error",
			giveErr: fmt.Errorf("wraps err: %w", &data.BaseError{
				Code:    data.ErrCodeFriendly,
				Message: "friendly error",
			}),
			wantResp: &data.Response{
				Code:    data.ErrCodeFriendly,
				Message: "wraps err: friendly error",
			},
		},
		{
			name: "coded-error",
			giveErr: fmt.Errorf("wraps err: %w", &codedTestError{
				code:    20001,
				message: "validation failed",
			}),
			wantResp: &data.Response{
				Code:    20001,
				Message: "wraps err: validation failed",
			},
		},
		{
			name: "coded-error-data",
			giveErr: &extraDataTestError{
				Cause: &codedTestError{
					code:    20001,
					message: "validation failed",
				},
				Data: map[string]string{"field": "tenant_id"},
			},
			wantResp: &data.Response{
				Code:    20001,
				Message: "validation failed",
				Data:    map[string]string{"field": "tenant_id"},
			},
		},
		{
			name:     "text-resp",
			giveResp: "test",
			wantResp: &data.Response{
				Data: "test",
			},
		},
		{
			name: "wrapper-resp",
			giveResp: &data.Response{
				Data: "test",
			},
			wantResp: &data.Response{
				Data: "test",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mMw := &core.MockMiddleware{}
			mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
				ctx := args.Get(0).(core.Context)
				ctx.SetOutput(tc.giveResp)
			}).Return(tc.giveErr)

			testMw := HttpRespReshapeMiddleware{
				opt: HttpRespReshapeOpt{DefaultErrCode: tc.giveOptCode},
				BaseMiddleware: BaseMiddleware{
					next: mMw,
				},
				respNewFunc: func() data.HttpResponse {
					return &data.Response{}
				},
			}
			c := core.NewContext()
			err := testMw.Handle(c)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Equal(t, tc.wantResp, c.Output())
		})
	}
}
