package middleware

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/shiningrush/droplet/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHttpInfoInjectorMiddleware(t *testing.T) {
	tests := []struct {
		name                   string
		header                 http.Header
		headerKeyRequestID     string
		wantTraceID            string
		wantRequestID          string
		wantOperator           string
		wantTenantID           string
		wantAppID              string
		wantLocale             string
		wantRemainingTimeoutMS string
	}{
		{
			name:               "prefer standard request id",
			headerKeyRequestID: "X-Request-ID",
			header: http.Header{
				"ofa-pass-trace-id":               []string{"trace-id"},
				"ofa-direct-request-id":           []string{"ofa-req-id"},
				"X-Request-ID":                    []string{"configured-req-id"},
				"ofa-pass-operator":               []string{"operator"},
				"ofa-pass-tenant-id":              []string{"tenant"},
				"ofa-pass-app-id":                 []string{"app"},
				"ofa-pass-locale":                 []string{"zh-CN"},
				"ofa-direct-remaining-timeout-ms": []string{"1000"},
			},
			wantTraceID:            "trace-id",
			wantRequestID:          "ofa-req-id",
			wantOperator:           "operator",
			wantTenantID:           "tenant",
			wantAppID:              "app",
			wantLocale:             "zh-CN",
			wantRemainingTimeoutMS: "1000",
		},
		{
			name:               "fallback to configured request id header",
			headerKeyRequestID: "X-Request-Id",
			header: http.Header{
				"X-Request-Id": []string{"configured-req-id"},
			},
			wantRequestID: "configured-req-id",
		},
		{
			name: "ignore old OFA request id header",
			header: http.Header{
				"OFA_DIRECT_REQUEST_ID": []string{"legacy-req-id"},
			},
		},
		{
			name: "generate missing ids",
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
				HeaderKeyRequestID: tc.headerKeyRequestID,
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
				if tc.wantTraceID != "" {
					assert.Equal(t, tc.wantTraceID, ctx.Get(KeyTraceID))
					assert.Equal(t, tc.wantTraceID, ctx.Context().Value(KeyTraceID))
				} else {
					assert.Regexp(t, regexp.MustCompile(`^[0-9a-f]{32}$`), ctx.Get(KeyTraceID))
					assert.Regexp(t, regexp.MustCompile(`^[0-9a-f]{32}$`), ctx.Context().Value(KeyTraceID))
				}
				if tc.wantRequestID != "" {
					assert.Equal(t, tc.wantRequestID, ctx.Get(KeyRequestID))
					assert.Equal(t, tc.wantRequestID, ctx.Context().Value(KeyRequestID))
				} else {
					assert.Regexp(t, regexp.MustCompile(`^req_[0-9]{8}_[0-9]{6}_[a-z2-7]{16}$`), ctx.Get(KeyRequestID))
					assert.Regexp(t, regexp.MustCompile(`^req_[0-9]{8}_[0-9]{6}_[a-z2-7]{16}$`), ctx.Context().Value(KeyRequestID))
				}
				assert.Equal(t, req.URL.Path, ctx.Path())
				assert.Nil(t, ctx.Request())
				assert.Equal(t, tc.wantOperator, ctx.GetString(KeyOperator))
				assert.Equal(t, tc.wantOperator, stringFromContext(ctx.Context(), KeyOperator))
				assert.Equal(t, tc.wantTenantID, ctx.GetString(KeyTenantID))
				assert.Equal(t, tc.wantTenantID, stringFromContext(ctx.Context(), KeyTenantID))
				assert.Equal(t, tc.wantAppID, ctx.GetString(KeyAppID))
				assert.Equal(t, tc.wantAppID, stringFromContext(ctx.Context(), KeyAppID))
				assert.Equal(t, tc.wantLocale, ctx.GetString(KeyLocale))
				assert.Equal(t, tc.wantLocale, stringFromContext(ctx.Context(), KeyLocale))
				assert.Equal(t, tc.wantRemainingTimeoutMS, ctx.GetString(KeyRemainingTimeoutMS))
				assert.Equal(t, tc.wantRemainingTimeoutMS, stringFromContext(ctx.Context(), KeyRemainingTimeoutMS))
			}).Return(nil)
			h.SetNext(mMw)

			_ = h.Handle(core.NewContext())
			assert.True(t, nextCalled)
		})
	}
}

func TestHttpInfoInjectorMiddlewarePrefersExistingContextValues(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{Path: "path"},
		Header: http.Header{
			"ofa-pass-trace-id":               []string{"header-trace"},
			"ofa-direct-request-id":           []string{"header-req"},
			"ofa-pass-operator":               []string{"header-operator"},
			"ofa-pass-tenant-id":              []string{"header-tenant"},
			"ofa-pass-app-id":                 []string{"header-app"},
			"ofa-pass-locale":                 []string{"en-US"},
			"ofa-direct-remaining-timeout-ms": []string{"1000"},
			"X-Request-ID":                    []string{"configured-req"},
		},
	}
	h := NewHttpInfoInjectorMiddleware(HttpInfoInjectorOption{
		HeaderKeyRequestID: "X-Request-ID",
		ReqFunc: func() *http.Request {
			return req
		},
	})

	rawCtx := context.Background()
	rawCtx = context.WithValue(rawCtx, KeyTraceID, "context-trace")
	rawCtx = context.WithValue(rawCtx, KeyRequestID, "context-req")
	rawCtx = context.WithValue(rawCtx, KeyOperator, "context-operator")
	rawCtx = context.WithValue(rawCtx, KeyTenantID, "context-tenant")
	rawCtx = context.WithValue(rawCtx, KeyAppID, "context-app")
	rawCtx = context.WithValue(rawCtx, KeyLocale, "zh-CN")
	rawCtx = context.WithValue(rawCtx, KeyRemainingTimeoutMS, "2000")

	nextCalled := false
	mMw := &core.MockMiddleware{}
	mMw.On("Handle", mock.Anything).Run(func(args mock.Arguments) {
		nextCalled = true
		ctx := args.Get(0).(core.Context)
		assert.Equal(t, "context-trace", ctx.GetString(KeyTraceID))
		assert.Equal(t, "context-trace", ctx.Context().Value(KeyTraceID))
		assert.Equal(t, "context-req", ctx.GetString(KeyRequestID))
		assert.Equal(t, "context-req", ctx.Context().Value(KeyRequestID))
		assert.Equal(t, "context-operator", ctx.GetString(KeyOperator))
		assert.Equal(t, "context-tenant", ctx.GetString(KeyTenantID))
		assert.Equal(t, "context-app", ctx.GetString(KeyAppID))
		assert.Equal(t, "zh-CN", ctx.GetString(KeyLocale))
		assert.Equal(t, "2000", ctx.GetString(KeyRemainingTimeoutMS))
	}).Return(nil)
	h.SetNext(mMw)

	c := core.NewContext()
	c.SetContext(rawCtx)
	_ = h.Handle(c)
	assert.True(t, nextCalled)
}

func stringFromContext(ctx context.Context, key string) string {
	value, _ := ctx.Value(key).(string)
	return value
}
