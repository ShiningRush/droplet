package codec

import (
	"bytes"
	"fmt"
	"github.com/shiningrush/droplet/data"
	"io"
	"io/ioutil"
	"net/http"
)

type MultipartForm struct {
}

func (j *MultipartForm) ContentType() []string {
	return []string{"multipart/form-data"}
}

func (j *MultipartForm) UnmarshalSearchMap(req *http.Request, ptr interface{}) (SearchMap, error) {
	// make request body could read multiple
	bs, err := data.CopyBody(req)
	if err != nil {
		return nil, err
	}

	reader, err := req.MultipartReader()
	if err != nil {
		return nil, fmt.Errorf("read form-data input from body failed: %s", err)
	}

	multiParts := SearchMap{}
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read next part from body failed: %s", err)
		}

		bs, err := ioutil.ReadAll(p)
		if err != nil {
			return nil, fmt.Errorf("read part body from body failed: %s", err)
		}
		if p.FileName() != "" {
			multiParts[fmt.Sprintf("%s%s", "_", p.FormName())] = []byte(p.FileName())
		}

		multiParts[p.FormName()] = bs
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
	return multiParts, nil
}

func (j *MultipartForm) Marshal(ptr interface{}) ([]byte, error) {
	panic("MultipartForm not support marshal")
}
