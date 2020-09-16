package data

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

func (r *Response) Set(code int, msg, reqId string, data interface{}) {
	r.Code = code
	r.Message = msg
	r.RequestID = reqId
	r.Data = data
}

type FileResponse struct {
	Name        string
	ContentType string
	Content     []byte
}

func (r *FileResponse) Get() (name, contentType string, content []byte) {
	return r.Name, r.ContentType, r.Content
}
