package data

import (
	"fmt"
	"net/http"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

func (r *Response) Set(code int, msg string, data interface{}) {
	r.Code = code
	r.Message = msg
	r.Data = data
}

func (r *Response) SetReqID(reqId string) {
	r.RequestID = reqId
}

type SpecCodeResponse struct {
	Response
	StatusCode int `json:"-"`
}

func (r *SpecCodeResponse) GetStatusCode() int {
	return r.StatusCode
}

type FileResponse struct {
	Name        string
	ContentType string
	Content     []byte
}

func (r *FileResponse) Get() (name, contentType string, content []byte) {
	return r.Name, r.ContentType, r.Content
}

type RawResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (rr *RawResponse) WriteRawResponse(rw http.ResponseWriter) error {
	for k, v := range rr.Header {
		rw.Header()[k] = v
	}
	rw.WriteHeader(rr.StatusCode)
	if _, err := rw.Write(rr.Body); err != nil {
		return fmt.Errorf("write body failed: %w", err)
	}
	return nil
}
