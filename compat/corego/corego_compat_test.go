//go:build compat_corego
// +build compat_corego

package corego_compat_check

import (
	"fmt"
	"testing"

	"github.com/dev-ofa/core-go/model/datax"
	dcore "github.com/shiningrush/droplet/core"
	ddata "github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/middleware"
)

type stubMiddleware struct {
	output interface{}
	err    error
}

func (m *stubMiddleware) SetNext(next dcore.Middleware) {}

func (m *stubMiddleware) Handle(ctx dcore.Context) error {
	ctx.SetOutput(m.output)
	return m.err
}

func TestDropletDataHelpersSupportRealCoreGoErrors(t *testing.T) {
	err := fmt.Errorf("wrap: %w", datax.NewValidationError("bad request", []datax.ValidateErrItem{
		{ParamName: "tenant_id", Reason: "missing"},
	}, nil))
	err = datax.WithErrorData(err, map[string]string{"field": "tenant_id"})

	if !ddata.IsErrCode(datax.ErrCodeValidate, err) {
		t.Fatalf("droplet IsErrCode should recognize core-go validation error")
	}
	if got := ddata.CodeOf(err); got != datax.ErrCodeValidate {
		t.Fatalf("droplet CodeOf = %d, want %d", got, datax.ErrCodeValidate)
	}

	data, ok := ddata.ErrorData(err).(map[string]string)
	if !ok {
		t.Fatalf("droplet ErrorData type = %T, want map[string]string", ddata.ErrorData(err))
	}
	if got := data["field"]; got != "tenant_id" {
		t.Fatalf("droplet ErrorData field = %q, want tenant_id", got)
	}
}

func TestDropletRespReshapeSupportsRealCoreGoErrors(t *testing.T) {
	coreErr := datax.WithErrorData(
		datax.NewValidationError("bad request", []datax.ValidateErrItem{
			{ParamName: "tenant_id", Reason: "missing"},
		}, nil),
		map[string]string{"field": "tenant_id"},
	)

	testMw := middleware.NewRespReshapeMiddleware(func() ddata.HttpResponse {
		return &ddata.Response{}
	}, middleware.HttpRespReshapeOpt{})
	testMw.SetNext(&stubMiddleware{err: coreErr})

	ctx := dcore.NewContext()
	ctx.Set(middleware.KeyRequestID, "req-1")

	if err := testMw.Handle(ctx); err != nil {
		t.Fatalf("reshape middleware returned err: %v", err)
	}

	resp, ok := ctx.Output().(*ddata.Response)
	if !ok {
		t.Fatalf("reshape output type = %T, want *data.Response", ctx.Output())
	}
	if resp.Code != datax.ErrCodeValidate {
		t.Fatalf("reshape code = %d, want %d", resp.Code, datax.ErrCodeValidate)
	}
	if resp.Message != "bad request" {
		t.Fatalf("reshape message = %q, want bad request", resp.Message)
	}
	if resp.RequestID != "req-1" {
		t.Fatalf("reshape request_id = %q, want req-1", resp.RequestID)
	}
	data, ok := resp.Data.(map[string]string)
	if !ok {
		t.Fatalf("reshape data type = %T, want map[string]string", resp.Data)
	}
	if got := data["field"]; got != "tenant_id" {
		t.Fatalf("reshape data field = %q, want tenant_id", got)
	}
}
