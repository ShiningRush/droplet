package middleware

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInputMiddleWare_InjectFieldFromBody(t *testing.T) {

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
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret, err := changeToFieldKind(tc.giveVal, tc.giveField)
			tc.wantErr(t, err)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}
