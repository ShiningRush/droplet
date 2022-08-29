package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	handler *ServerTestHandler
}

func (sts *ServerTestSuite) TestJsonResponse() {

	input := JsonResponseParam{
		IntBody:    10,
		StrBody:    "hello",
		Path:       "path",
		Query:      "query1",
		UrlDefault: "query2",
	}

	reqBody, _ := json.Marshal(input)
	bd, err := http.Post(sts.handler.Addr(), "application/json", bytes.NewBuffer(reqBody))
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), http.StatusOK, bd.StatusCode)

	bdBytes, err := ioutil.ReadAll(bd.Body)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), "{}", bdBytes)
}

type ServerTestHandler struct {
}

func (h *ServerTestHandler) Addr() string {
	return "127.0.0.1:19601"
}

type JsonResponseParam struct {
	IntBody    int    `json:"intBody"`
	StrBody    string `json:"strBody"`
	Path       string `auth_read:"path,path"`
	Query      string `auth_read:"path,query"`
	UrlDefault string `auth_read:"default"`
}

func (h *ServerTestHandler) JsonResponse(ctx core.Context) (interface{}, error) {
	return ctx.Input(), nil
}

type SpecResponseOption struct {
	Type string
}

func (h *ServerTestHandler) SpecResponse(ctx core.Context) (interface{}, error) {
	opt := ctx.Input().(*SpecResponseOption)

	switch opt.Type {
	case "err-text":
		return nil, errors.New("text error")
	case "err-not-found":
		return nil, data.ErrNotFound
	case "spec":
		return data.SpecCodeResponse{
			StatusCode: 400,
			Response: data.Response{
				Code:    20000,
				Message: "spec-msg",
				Data: map[string]interface{}{
					"demoKey": "demoVal",
				},
			},
		}, nil
	}
	return nil, nil
}
