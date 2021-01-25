package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CopyBody(req *http.Request) ([]byte, error) {
	// make request body could read multiple
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}
	if err := req.Body.Close(); err != nil {
		return nil, fmt.Errorf("close body failed: %w", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}
