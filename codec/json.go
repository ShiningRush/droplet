package codec

import (
	"encoding/json"
	"net/http"

	"github.com/shiningrush/droplet/data"
)

type Json struct {
}

func (j *Json) ContentType() []string {
	return []string{"application/json"}
}

func (j *Json) Unmarshal(req *http.Request, ptr interface{}) error {
	bs, err := data.CopyBody(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, ptr)
}

func (j *Json) Marshal(ptr interface{}) ([]byte, error) {
	return json.Marshal(ptr)
}
