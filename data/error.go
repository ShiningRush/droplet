package data

import (
	"fmt"
)

const (
	ErrCodeInternal = 10000
	ErrCodeNotFound = 10001
	ErrCodeConflict = 10002
	ErrCodeFriendly = 10003
	ErrCodeValidate = 10004
	ErrCodeFormat   = 10005
)

var (
	ErrNotFound = &BaseError{Code: ErrCodeNotFound, Message: "data not found"}
	ErrConflict = &BaseError{Code: ErrCodeConflict, Message: "data is existed or has be updated"}
)

type BaseError struct {
	Code      int
	Message   string
	Data      interface{}
	SourceSrv string
}

func (e *BaseError) Error() string {
	if e.SourceSrv != "" {
		return fmt.Sprintf("call: %s failed, code: %d, msg: %s", e.SourceSrv, e.Code, e.Message)
	}
	return e.Message
}

func (e *BaseError) Is(err error) bool {
	if err == nil {
		return false
	}
	switch t := err.(type) {
	case *CallSrvError:
		return t.Code == e.Code
	case *BaseError:
		return t.Code == e.Code
	}
	return false
}

type CallSrvError struct {
	SrvResp *Response
	BaseError
}

func NewNotFoundError(msg string) error {
	if msg == "" {
		return ErrNotFound
	}
	return &BaseError{
		Code:    100,
		Message: msg,
	}
}

func NewConflictError(msg string) error {
	if msg == "" {
		return ErrConflict
	}
	return &BaseError{
		Code:    101,
		Message: msg,
	}
}

func NewInternalError(msg string) error {
	if msg == "" {
		msg = "internal server error"
	}
	return &BaseError{
		Code:    ErrCodeInternal,
		Message: msg,
	}
}

func NewFriendlyError(msg string) error {
	return &BaseError{
		Code:    ErrCodeFriendly,
		Message: msg,
	}
}

func NewFormatError(msg string) error {
	return &BaseError{
		Code:    ErrCodeFormat,
		Message: msg,
	}
}

type ValidateError struct {
	BaseError
}

func NewValidateError(msg string, items []ValidateErrItem) error {
	if msg == "" {
		msg = "parameter validate failed"
	}
	return &BaseError{
		Code:    ErrCodeValidate,
		Message: msg,
		Data:    items,
	}
}

type ValidateErrItem struct {
	ParamName string      `json:"paramName"`
	Reason    string      `json:"reason"`
	Detail    interface{} `json:"detail"`
}
