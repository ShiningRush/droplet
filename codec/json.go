package codec

import (
	"encoding/json"
	"github.com/shiningrush/droplet/data"
	"net/http"
)

type Json struct {
}

func (j *Json) ContentType() string {
	return "application/json"
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
