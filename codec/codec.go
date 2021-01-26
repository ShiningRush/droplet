package codec

import "net/http"

type Interface interface {
	ContentType() []string
	Marshal(interface{}) ([]byte, error)
}

type Direct interface {
	Interface
	Unmarshal(*http.Request, interface{}) error
}

type SearchMap map[string][]byte
type Search interface {
	Interface
	UnmarshalSearchMap(*http.Request) (SearchMap, error)
}
