package wrapper

import (
	"reflect"
	"testing"
)

func TestInputTypeof(t *testing.T) {
	type testInputType struct {
		A string
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want WrapOptBase
	}{
		{
			name: "struct",
			args: args{
				v: testInputType{},
			},
			want: WrapOptBase{
				inputType: reflect.TypeOf(testInputType{}),
			},
		},
		{
			name: "ptr",
			args: args{
				v: &testInputType{},
			},
			want: WrapOptBase{
				inputType: reflect.TypeOf(testInputType{}),
			},
		},
		{
			name: "nil",
			args: args{
				v: nil,
			},
			want: WrapOptBase{},
		},
		{
			name: "nil Type",
			args: args{
				v: (*testInputType)(nil),
			},
			want: WrapOptBase{
				inputType: reflect.TypeOf(testInputType{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotOpt = InputTypeOf(tt.args.v)
			got := WrapOptBase{}
			gotOpt(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputTypeOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
