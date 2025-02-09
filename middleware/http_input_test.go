package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FixedPathFunc(key string) string {
	return key
}

type ValidateInput struct {
	RequiredStr string `validate:"required"`

	validateErr error
}

func (input *ValidateInput) Initial(ctx core.Context) error {
	return input.validateErr
}

func TestInputMiddleWare_inputValidate(t *testing.T) {
	tests := []struct {
		name      string
		giveMw    *HttpInputMiddleware
		giveInput interface{}
		wantErr   error
	}{
		{
			name: "normal",
			giveMw: NewHttpInputMiddleWare(HttpInputOption{
				ValidateErrCode: data.ErrCodeValidate,
			}),
			giveInput: &ValidateInput{RequiredStr: "test"},
			wantErr:   nil,
		},
		{
			name: "validate failed",
			giveMw: NewHttpInputMiddleWare(HttpInputOption{
				ValidateErrCode: data.ErrCodeValidate,
			}),
			giveInput: &ValidateInput{},
			wantErr: &data.BaseError{
				Code:    data.ErrCodeValidate,
				Message: "input validate failed: Key: 'ValidateInput.RequiredStr' Error:Field validation for 'RequiredStr' failed on the 'required' tag",
			},
		},
		{
			name: "initial failed",
			giveMw: NewHttpInputMiddleWare(HttpInputOption{
				ValidateErrCode: data.ErrCodeInternal,
			}),
			giveInput: &ValidateInput{validateErr: fmt.Errorf("some err")},
			wantErr: &data.BaseError{
				Code:    data.ErrCodeInternal,
				Message: "input initial failed: some err",
			},
		},
		{
			name: "custom err",
			giveMw: NewHttpInputMiddleWare(HttpInputOption{
				ValidateErrCode: data.ErrCodeInternal,
			}),
			giveInput: &ValidateInput{validateErr: data.NewConflictError("err")},
			wantErr: &data.BaseError{
				Code:    data.ErrCodeConflict,
				Message: "input initial failed: err",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.giveMw.inputValidate(nil, tc.giveInput)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

type TestInput struct {
	QueryString   string  `auto_read:"query_str,query"`
	QueryEmptyPtr *string `auto_read:"empty_str,query"`
	HeaderInt     int     `auto_read:"header-int,header"`
	DefaultIntPtr *int    `auto_read:"query_int"`
	PathStrPtr    *string `auto_read:"path_str,path"`
	CookieStrPtr  *string `auto_read:"cookie_str,cookie"`
	MixedStrPtr   *string `auto_read:"mixed, query|header|cookie"`
	MixedStr      string  `auto_read:"mixed_str, query|header|cookie"`
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
					URL:    &url.URL{RawQuery: "query_str=query_string&test=2&mixed=query&mixed_str=query_str"},
					Method: http.MethodPost,
					Header: map[string][]string{
						"Header-Int": {"10"},
						"Cookie":     {"cookie_str=c_str;mixed=cookie"},
						"Mixed":      {"header"},
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
				CookieStrPtr:  strPtr("c_str"),
				Body:          []byte("all body"),
				MixedStrPtr:   strPtr("query"),
				MixedStr:      "query_str",
			},
			wantErr: require.NoError,
		},
		{
			name: "mixed-header",
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
						"Cookie":     {"cookie_str=c_str;mixed=cookie"},
						"Mixed":      {"header"},
						"Mixed_str":  {"header_str"},
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
				CookieStrPtr:  strPtr("c_str"),
				Body:          []byte("all body"),
				MixedStrPtr:   strPtr("header"),
				MixedStr:      "header_str",
			},
			wantErr: require.NoError,
		},
		{
			name: "mixed-cookie",
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
						"Cookie":     {"cookie_str=c_str;mixed=cookie;mixed_str=cookie_str"},
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
				CookieStrPtr:  strPtr("c_str"),
				Body:          []byte("all body"),
				MixedStrPtr:   strPtr("cookie"),
				MixedStr:      "cookie_str",
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
