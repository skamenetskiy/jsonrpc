package jsonrpc

import (
	"context"
	"encoding/json"
	"net/http"
)

type Request interface {
	Bind(v interface{}) error
	Context() context.Context
	Header(name string) string
	Request() *http.Request
}

type requestBody struct {
	Method    string          `json:"method"`
	Data      json.RawMessage `json:"data"`
	Signature []byte          `json:"signature"`
}

type request struct {
	req  *http.Request
	body *requestBody
}

func (req *request) Bind(v interface{}) error {
	return json.Unmarshal(req.body.Data, v)
}

func (req *request) Context() context.Context {
	return req.Request().Context()
}

func (req *request) Header(name string) string {
	return req.Request().Header.Get(name)
}

func (req *request) Request() *http.Request {
	return req.req
}
