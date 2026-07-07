package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/shiningrush/droplet/core"
)

const ofaDirectRequestIDHeader = "OFA_DIRECT_REQUEST_ID"

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

func requestIDFromHeaders(header http.Header, configuredHeader string) string {
	for _, key := range requestIDHeaderCandidates(configuredHeader) {
		if reqID := header.Get(key); reqID != "" {
			return reqID
		}
	}
	return ""
}

func requestIDHeaderCandidates(configuredHeader string) []string {
	candidates := []string{ofaDirectRequestIDHeader}
	if configuredHeader != "" && !strings.EqualFold(configuredHeader, ofaDirectRequestIDHeader) {
		candidates = append(candidates, configuredHeader)
	}
	return candidates
}

func (mw *HttpInfoInjectorMiddleware) Handle(ctx core.Context) error {
	req := mw.opt.ReqFunc()

	ctx.Set(KeyHttpRequest, req)
	if ctx.Get(KeyRequestID) == nil {
		reqId := requestIDFromHeaders(req.Header, mw.opt.HeaderKeyRequestID)
		if reqId == "" {
			reqId = fmt.Sprintf("%s%v", time.Now().UTC().Format("20060102150405"), rand.Intn(100000))
		}
		ctx.Set(KeyRequestID, reqId)
	}
	ctx.SetPath(req.URL.Path)

	return mw.BaseMiddleware.Handle(ctx)
}
