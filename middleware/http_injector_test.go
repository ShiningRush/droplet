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
		name          string
		header        http.Header
		headerKey     string
		wantRequestID string
	}{
		{
			name:      "prefer OFA request id header",
			headerKey: "X-Request-ID",
			header: http.Header{
				"OFA_DIRECT_REQUEST_ID": []string{"ofa-req-id"},
				"X-Request-ID":          []string{"legacy-req-id"},
			},
			wantRequestID: "ofa-req-id",
		},
		{
			name:      "fallback to configured header",
			headerKey: "X-Request-Id",
			header: http.Header{
				"X-Request-Id": []string{"reqId"},
			},
			wantRequestID: "reqId",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			header := http.Header{}
			for key, values := range tc.header {
				for _, value := range values {
					header.Add(key, value)
				}
			}
			req := &http.Request{
				URL: &url.URL{
					Path: "path",
				},
				Header: header,
			}
			h := NewHttpInfoInjectorMiddleware(HttpInfoInjectorOption{
				HeaderKeyRequestID: tc.headerKey,
				ReqFunc: func() *http.Request {
					return req
				},
			})

			nextCalled := false
			mMw := &core.MockMiddleware{}
			mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
				nextCalled = true
				ctx := args.Get(0).(core.Context)
				assert.Same(t, req, ctx.Get(KeyHttpRequest))
				assert.Equal(t, tc.wantRequestID, ctx.Get(KeyRequestID))
				assert.Equal(t, req.URL.Path, ctx.Path())
				assert.Nil(t, ctx.Request())
			}).Return(nil)
			h.SetNext(mMw)

			_ = h.Handle(core.NewContext())
			assert.True(t, nextCalled)
		})
	}
}
