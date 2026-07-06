package wrapper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shiningrush/droplet/core"
)

const (
	serverTimingHeader = "Server-Timing"
	appTimingMetric    = "app"
)

func TestHandleHttpInPipelineSetsServerTimingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/books", nil)
	req.Header.Set("X-Request-ID", "req-1")
	recorder := httptest.NewRecorder()

	HandleHttpInPipeline(HandleHttpInPipelineInput{
		Req:            req,
		RespWriter:     NewResponseWriter(recorder),
		PathParamsFunc: func(string) string { return "" },
		Handler: func(ctx core.Context) (interface{}, error) {
			return map[string]string{"status": "ok"}, nil
		},
	})

	headerValue := recorder.Header().Get(serverTimingHeader)
	if headerValue == "" {
		t.Fatal("Server-Timing header should be set")
	}
	if !strings.Contains(headerValue, appTimingMetric+";dur=") {
		t.Fatalf("Server-Timing header = %q, want %q metric", headerValue, appTimingMetric)
	}
}
