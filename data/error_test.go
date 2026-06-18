package data

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type codedTestError struct {
	code    int
	message string
	cause   error
}

func (e *codedTestError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *codedTestError) Unwrap() error {
	return e.cause
}

func (e *codedTestError) ErrorCode() int {
	return e.code
}

type extraDataTestError struct {
	Cause error
	Data  interface{}
}

func (e *extraDataTestError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "error with extra data"
}

func (e *extraDataTestError) Unwrap() error {
	return e.Cause
}

func TestIsErrCodeSupportsCodedError(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &codedTestError{
		code:    20001,
		message: "validation failed",
	})

	assert.True(t, IsErrCode(20001, err))
	assert.False(t, IsErrCode(ErrCodeInternal, errors.New("plain")))
}

func TestCodeOf(t *testing.T) {
	assert.Equal(t, 0, CodeOf(nil))
	assert.Equal(t, ErrCodeInternal, CodeOf(errors.New("plain")))
	assert.Equal(t, 20002, CodeOf(fmt.Errorf("wrap: %w", &codedTestError{
		code:    20002,
		message: "resource not found",
	})))
}

func TestErrorData(t *testing.T) {
	assert.Equal(t, "base data", ErrorData(fmt.Errorf("wrap: %w", &BaseError{
		Code:    ErrCodeFriendly,
		Message: "friendly error",
		Data:    "base data",
	})))

	assert.Equal(t, "wrapper data", ErrorData(&ErrWrapper{
		Code: 20000,
		Msg:  "wrapper error",
		Data: "wrapper data",
	}))

	err := &extraDataTestError{
		Cause: &codedTestError{
			code:    20001,
			message: "validation failed",
		},
		Data: map[string]string{"field": "tenant_id"},
	}
	assert.Equal(t, map[string]string{"field": "tenant_id"}, ErrorData(fmt.Errorf("wrap: %w", err)))
}
