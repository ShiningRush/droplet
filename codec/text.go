package codec

import (
        "github.com/shiningrush/droplet/data"
        "net/http"
)

type Text struct {
}

func (j *Text) ContentType() []string {
        return []string{"text/plain", "text/html"}
}

func (j *Text) Unmarshal(req *http.Request, ptr interface{}) error {
        bs, err := data.CopyBody(req)
        if err != nil {
                return err
        }

        ptr = bs

        return nil
}

func (j *Text) Marshal(ptr interface{}) ([]byte, error) {
        panic("MultipartForm not support marshal")
}
