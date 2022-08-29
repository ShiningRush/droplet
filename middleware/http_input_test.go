package middleware

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		wantErr   error
	}{
		{
			caseDesc:  "bool",
			giveVal:   "true",
			giveField: reflect.TypeOf(true),
			wantRet:   true,
		},
		{
			caseDesc:  "boolPtr",
			giveVal:   "true",
			giveField: reflect.TypeOf(boolPtr(true)),
			wantRet:   boolPtr(true),
		},
		{
			caseDesc:  "empty boolPtr",
			giveVal:   "",
			giveField: reflect.TypeOf(boolPtr(true)),
			wantRet:   nil,
		},
		{
			caseDesc:  "init boolPtr",
			giveVal:   "true",
			giveField: reflect.TypeOf(initBoolPtr),
			wantRet:   boolPtr(true),
		},
		{
			caseDesc:  "int",
			giveVal:   "123",
			giveField: reflect.TypeOf(123),
			wantRet:   123,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret, err := changeToFieldKind(tc.giveVal, tc.giveField)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}
