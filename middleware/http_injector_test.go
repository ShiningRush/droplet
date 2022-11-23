package middleware

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/shiningrush/droplet/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHttpInfoInjectorMiddleware(t *testing.T) {
	tests := []struct {
		giveOpt HttpInfoInjectorOption
		giveCtx core.Context
	}{
		{
			giveCtx: core.NewContext(nil),
			giveOpt: HttpInfoInjectorOption{
				ReqFunc: func() *http.Request {
					return &http.Request{
						URL: &url.URL{
							Path: "path",
						},
						Header: http.Header{
							"X-Request-ID": []string{"reqId"},
						},
					}
				},
			},
		},
	}

	for _, tc := range tests {
		h := NewHttpInfoInjectorMiddleware(tc.giveOpt)

		nextCalled := false
		mMw := &core.MockMiddleware{}
		mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
			nextCalled = true
			ctx := args.Get(0).(core.Context)
			req := ctx.Get(KeyHttpRequest)
			assert.Equal(t, tc.giveOpt.ReqFunc(), req)
			reqId := ctx.Get(KeyRequestID)
			assert.Equal(t, tc.giveOpt.ReqFunc().Header.Get("X-Request-ID"), reqId)
			assert.Equal(t, ctx.Path(), tc.giveOpt.ReqFunc().URL.Path)
			assert.Equal(t, tc.giveCtx.Request(), ctx.Request())
		}).Return(nil)
		h.SetNext(mMw)

		_ = h.Handle(tc.giveCtx)
		assert.True(t, nextCalled)
	}
}
