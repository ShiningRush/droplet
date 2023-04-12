package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FixedPathFunc(key string) string {
	return key
}

type TestInput struct {
	QueryString   string  `auto_read:"query_str,query"`
	HeaderInt     int     `auto_read:"header-int,header"`
	DefaultIntPtr *int    `auto_read:"query_int"`
	PathStrPtr    *string `auto_read:"path_str,path"`
	Body          []byte  `auto_read:"@body"`
}

func strPtr(str string) *string {
	return &str
}

func TestInputMiddleWare_injectFieldFromUrlAndMap(t *testing.T) {
	tests := []struct {
		name    string
		giveMw  *HttpInputMiddleware
		givePtr interface{}
		wantPtr interface{}
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "normal",
			giveMw: &HttpInputMiddleware{
				opt: HttpInputOption{
					PathParamsFunc: FixedPathFunc,
					IsReadFromBody: true,
				},
				req: &http.Request{
					URL:    &url.URL{RawQuery: "query_str=query_string&test=2"},
					Method: http.MethodPost,
					Header: map[string][]string{
						"Header-Int": {"10"},
					},
					Body: io.NopCloser(bytes.NewBufferString("all body")),
				},
				searchMap: nil,
			},
			givePtr: &TestInput{},
			wantPtr: &TestInput{
				QueryString:   "query_string",
				HeaderInt:     10,
				DefaultIntPtr: nil,
				PathStrPtr:    strPtr("path_str"),
				Body:          []byte("all body"),
			},
			wantErr: require.NoError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.giveMw.injectFieldFromUrlAndMap(tc.givePtr)
			tc.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantPtr, tc.givePtr)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func Test_changeToFieldKind(t *testing.T) {
	var initBoolPtr *bool
	tests := []struct {
		caseDesc  string
		giveVal   string
		giveField reflect.Type
		wantRet   interface{}
		wantErr   require.ErrorAssertionFunc
	}{
		{
			caseDesc:  "bool",
			giveVal:   "true",
			giveField: reflect.TypeOf(true),
			wantRet:   true,
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "boolPtr",
			giveVal:   "true",
			giveField: reflect.TypeOf(boolPtr(true)),
			wantRet:   boolPtr(true),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "empty boolPtr",
			giveVal:   "",
			giveField: reflect.TypeOf(boolPtr(true)),
			wantRet:   nil,
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "init boolPtr",
			giveVal:   "true",
			giveField: reflect.TypeOf(initBoolPtr),
			wantRet:   boolPtr(true),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "int",
			giveVal:   "123",
			giveField: reflect.TypeOf(123),
			wantRet:   123,
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "int64",
			giveVal:   "123",
			giveField: reflect.TypeOf(int64(123)),
			wantRet:   int64(123),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "int64-empty",
			giveVal:   "",
			giveField: reflect.TypeOf(int64(123)),
			wantRet:   int64(0),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "uint",
			giveVal:   "123",
			giveField: reflect.TypeOf(uint(123)),
			wantRet:   uint(123),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "uint-empty",
			giveVal:   "",
			giveField: reflect.TypeOf(uint(123)),
			wantRet:   uint(0),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "uint64",
			giveVal:   "123",
			giveField: reflect.TypeOf(uint64(123)),
			wantRet:   uint64(123),
			wantErr:   require.NoError,
		},
		{
			caseDesc:  "uint64-empty",
			giveVal:   "",
			giveField: reflect.TypeOf(uint64(123)),
			wantRet:   uint64(0),
			wantErr:   require.NoError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret, err := changeToFieldKind(tc.giveVal, tc.giveField)
			tc.wantErr(t, err)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}
