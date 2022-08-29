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
	// According to document, we does need to close request body:
	// For server requests, the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}
