package data

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type PagerInfo interface {
	GetPageInfo() (pageSize, pageNumber int, pageToken string)
	SetPageNumber(pageNumber int)
}

type SortInfo interface {
	GetSortInfo() []*SortPair
}

type SortAble struct {
	// format: "field[ desc], field2[ desc]"
	OrderBy string `json:"order_by" form:"order_by" auto_read:"order_by"`
}

func (s *SortAble) GetSortInfo() (paris []*SortPair) {
	if s.OrderBy == "" {
		return
	}

	sortFields := strings.Split(s.OrderBy, ",")
	for _, aField := range sortFields {
		nameAndDesc := strings.Split(strings.TrimSpace(aField), " ")
		p := &SortPair{Field: nameAndDesc[0]}
		if len(nameAndDesc) > 1 && nameAndDesc[1] == "desc" {
			p.IsDescending = true
		}
		paris = append(paris, p)
	}

	return
}

type SortPair struct {
	Field        string
	IsDescending bool
}

// BuildNextPageToken It is google web API, please refer https://cloud.google.com/apis/design/design_patterns#list_pagination
func BuildNextPageToken(pager PagerInfo) (string, error) {
	_, pn, _ := pager.GetPageInfo()
	pager.SetPageNumber(pn + 1)
	bs, err := json.Marshal(pager)
	if err != nil {
		return "", fmt.Errorf("to token faild when marshal pager: %s", err)
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

func RecoverPager(pager PagerInfo) (bool, error) {
	_, _, token := pager.GetPageInfo()
	if token == "" {
		return false, nil
	}

	bs, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, fmt.Errorf("recover pager failed: %w", err)
	}
	err = json.Unmarshal(bs, pager)
	if err != nil {
		return false, fmt.Errorf("unmarshal pager failed: %w", err)
	}

	return true, nil
}
