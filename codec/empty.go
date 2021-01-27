package codec

import (
	"net/http"
)

type Empty struct {
}

func (j *Empty) ContentType() []string {
	return []string{"text/plain", "text/html"}
}

func (j *Empty) Unmarshal(req *http.Request, ptr interface{}) error {
	return nil
}

func (j *Empty) Marshal(ptr interface{}) ([]byte, error) {
	panic("Empty not support marshal")
}
