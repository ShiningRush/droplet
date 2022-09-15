package data

import (
	"fmt"
	"io"
	"net/http"
)

type HttpResponse interface {
	Set(code int, msg string, data interface{})
	SetReqID(reqId string)
}

type HttpFileResponse interface {
	Get() *FileResponse
}

type SpecCodeHttpResponse interface {
	GetStatusCode() int
	HttpResponse
}

type ResponseWriter interface {
	SetHeader(key, val string)
	GetHeader(key string) string
	GetHeaderValues(key string) []string
	DelHeader(key string)

	Write([]byte) (int, error)
	WriteHeader(statusCode int)

	// StdHttpWriter return the http.ResponseWriter, if wrapped framework is not compatible(such as fasthttp)
	// it will return nil
	StdHttpWriter() http.ResponseWriter
}

type RawHttpResponse interface {
	WriteRawResponse(writer ResponseWriter) error
}

// interface validate
var (
	_ HttpResponse         = (*Response)(nil)
	_ SpecCodeHttpResponse = (*SpecCodeResponse)(nil)
	_ HttpFileResponse     = (*FileResponse)(nil)
	_ RawHttpResponse      = (*RawResponse)(nil)
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
	Name          string
	ContentType   string
	Content       []byte
	ContentReader io.ReadCloser
	StatusCode    int
}

func (r *FileResponse) Get() *FileResponse {
	return r
}

type RawResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	BodyReader io.ReadCloser
}

func (rr *RawResponse) WriteRawResponse(rw ResponseWriter) error {
	for k, v := range rr.Header {
		rw.SetHeader(k, v[0])
	}
	rw.WriteHeader(rr.StatusCode)
	if rr.BodyReader != nil {
		defer rr.BodyReader.Close()
		if _, err := io.Copy(rw, rr.BodyReader); err != nil {
			return fmt.Errorf("copy body failed: %w", err)
		}
		return nil
	}

	if _, err := rw.Write(rr.Body); err != nil {
		return fmt.Errorf("write body failed: %w", err)
	}
	return nil
}
