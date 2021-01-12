package middleware

import (
	"github.com/shiningrush/droplet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

func TestNewHttpInfoInjectorMiddleware(t *testing.T) {
	tests := []struct {
		giveOpt HttpInfoInjectorOption
		giveCtx droplet.Context
	}{
		{
			giveCtx: droplet.NewContext(),
			giveOpt: HttpInfoInjectorOption{
				ReqFunc: func() *http.Request {
					return &http.Request{
						URL: &url.URL{
							Path: "path",
						},
						Header: http.Header{
							droplet.Option.HeaderKeyRequestID: []string{"reqId"},
						},
					}
				},
			},
		},
	}

	for _, tc := range tests {
		h := NewHttpInfoInjectorMiddleware(tc.giveOpt)

		nextCalled := false
		mMw := &droplet.MockMiddleware{}
		mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
			nextCalled = true
			ctx := args.Get(0).(droplet.Context)
			req := ctx.Get(KeyHttpRequest)
			assert.Equal(t, tc.giveOpt.ReqFunc(), req)
			reqId := ctx.Get(KeyRequestID)
			assert.Equal(t, tc.giveOpt.ReqFunc().Header.Get(droplet.Option.HeaderKeyRequestID), reqId)
			assert.Equal(t, ctx.Path(), tc.giveOpt.ReqFunc().URL.Path)
		}).Return(nil)
		h.SetNext(mMw)

		_ = h.Handle(tc.giveCtx)
		assert.True(t, nextCalled)
	}
}
