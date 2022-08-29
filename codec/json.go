package codec

import (
	"encoding/json"
	"fmt"
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

	if err := json.Unmarshal(bs, ptr); err != nil {
		return fmt.Errorf("json codec unmarshal failed: %w", err)
	}
	return nil
}

func (j *Json) Marshal(ptr interface{}) ([]byte, error) {
	bs, err := json.Marshal(ptr)
	if err != nil {
		return nil, fmt.Errorf("json codec marshal failed: %w", err)
	}
	return bs, nil
}
