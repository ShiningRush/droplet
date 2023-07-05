package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortAble_GetSortInfo(t *testing.T) {
	tests := []struct {
		name string
		give *SortAble
		want []*SortPair
	}{
		{
			name: "normal",
			give: &SortAble{OrderBy: "field1, field2 desc"},
			want: []*SortPair{
				{Field: "field1"},
				{Field: "field2", IsDescending: true},
			},
		},
		{
			name: "empty",
			give: &SortAble{OrderBy: ""},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := tt.give.GetSortInfo()
			assert.Equal(t, tt.want, sp)
		})
	}
}
