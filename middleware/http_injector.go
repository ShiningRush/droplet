package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shiningrush/droplet/core"
)

type HttpInfoInjectorMiddleware struct {
	BaseMiddleware
	opt HttpInfoInjectorOption
}

type HttpInfoInjectorOption struct {
	ReqFunc            func() *http.Request
	HeaderKeyRequestID string
}

func NewHttpInfoInjectorMiddleware(opt HttpInfoInjectorOption) *HttpInfoInjectorMiddleware {
	return &HttpInfoInjectorMiddleware{opt: opt}
}

func traceIDFromHeader(header http.Header) string {
	return header.Get(KeyTraceID)
}

func requestIDFromHeader(header http.Header) string {
	return header.Get(KeyRequestID)
}

func configuredRequestIDFromHeader(header http.Header, configuredHeader string) string {
	if configuredHeader == "" || strings.EqualFold(configuredHeader, KeyRequestID) {
		return ""
	}
	return header.Get(configuredHeader)
}

func generateTraceID() (string, error) {
	bs := make([]byte, 16)
	if _, err := rand.Read(bs); err != nil {
		return "", err
	}
	return hex.EncodeToString(bs), nil
}

func generateRequestID() (string, error) {
	bs := make([]byte, 10)
	if _, err := rand.Read(bs); err != nil {
		return "", err
	}
	suffix := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bs))
	return "req_" + time.Now().UTC().Format("20060102_150405") + "_" + suffix, nil
}

func contextValue(ctx core.Context, key string) string {
	if value := ctx.GetString(key); value != "" {
		return value
	}
	value, _ := ctx.Context().Value(key).(string)
	return value
}

func contextValueFromHeader(ctx core.Context, header http.Header, key string) string {
	if value := contextValue(ctx, key); value != "" {
		return value
	}
	return header.Get(key)
}

func setContextValue(ctx core.Context, key string, value string) {
	if value == "" {
		return
	}
	ctx.Set(key, value)
	ctx.SetContext(context.WithValue(ctx.Context(), key, value))
}

func (mw *HttpInfoInjectorMiddleware) Handle(ctx core.Context) error {
	req := mw.opt.ReqFunc()

	ctx.Set(KeyHttpRequest, req)
	traceID := contextValue(ctx, KeyTraceID)
	if traceID == "" {
		traceID = traceIDFromHeader(req.Header)
		if traceID == "" {
			var err error
			traceID, err = generateTraceID()
			if err != nil {
				return fmt.Errorf("generate trace id failed: %w", err)
			}
		}
	}
	ctx.Set(KeyTraceID, traceID)
	reqID := contextValue(ctx, KeyRequestID)
	if reqID == "" {
		reqID = requestIDFromHeader(req.Header)
		if reqID == "" {
			reqID = configuredRequestIDFromHeader(req.Header, mw.opt.HeaderKeyRequestID)
			if reqID == "" {
				var err error
				reqID, err = generateRequestID()
				if err != nil {
					return fmt.Errorf("generate request id failed: %w", err)
				}
			}
		}
	}
	ctx.Set(KeyRequestID, reqID)
	ctx.SetContext(context.WithValue(ctx.Context(), KeyTraceID, traceID))
	ctx.SetContext(context.WithValue(ctx.Context(), KeyRequestID, reqID))
	setContextValue(ctx, KeyOperator, contextValueFromHeader(ctx, req.Header, KeyOperator))
	setContextValue(ctx, KeyTenantID, contextValueFromHeader(ctx, req.Header, KeyTenantID))
	setContextValue(ctx, KeyAppID, contextValueFromHeader(ctx, req.Header, KeyAppID))
	setContextValue(ctx, KeyLocale, contextValueFromHeader(ctx, req.Header, KeyLocale))
	setContextValue(ctx, KeyRemainingTimeoutMS, contextValueFromHeader(ctx, req.Header, KeyRemainingTimeoutMS))
	ctx.SetPath(req.URL.Path)

	return mw.BaseMiddleware.Handle(ctx)
}
