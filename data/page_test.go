package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubPager struct {
	Filter string
	Pager
	SortAble
}

func TestBuildAndRecoverPager(t *testing.T) {
	s := &StubPager{}
	s.PageSize = 10
	s.PageNumber = 20
	s.Filter = "test filter"

	token, _ := BuildNextPageToken(s)
	newStub := &StubPager{}
	newStub.PageToken = token

	RecoverPager(newStub)
	assert.Equal(t, 10, newStub.PageSize)
	assert.Equal(t, 21, newStub.PageNumber)
	assert.Equal(t, "test filter", newStub.Filter)
}

func TestSortPair(t *testing.T) {
	s := &StubPager{}
	s.PageSize = 10
	s.PageNumber = 20
	s.Filter = "test filter"
	s.OrderBy = "  f1 desc  , f2  , f3 desc  ,    f4"

	sorts := s.GetSortInfo()
	assert.Equal(t, 4, len(sorts))
	for i := range sorts {
		if i == 0 {
			assert.Equal(t, "f1", sorts[i].Field)
			assert.Equal(t, true, sorts[i].IsDescending)
		}
		if i == 1 {
			assert.Equal(t, "f2", sorts[i].Field)
			assert.Equal(t, false, sorts[i].IsDescending)
		}
		if i == 2 {
			assert.Equal(t, "f3", sorts[i].Field)
			assert.Equal(t, true, sorts[i].IsDescending)
		}
		if i == 3 {
			assert.Equal(t, "f4", sorts[i].Field)
			assert.Equal(t, false, sorts[i].IsDescending)
		}
	}
}
