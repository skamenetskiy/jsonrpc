package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewResponse(status Status, data interface{}, err error) Response {
	return &response{
		status: status,
		data:   data,
		err:    err,
	}
}

func OK(data interface{}) Response {
	return NewResponse(StatusOK, data, nil)
}

func Error(err error) Response {
	return NewResponse(StatusError, nil, err)
}

func Errorf(format string, v ...interface{}) Response {
	return Error(fmt.Errorf(format, v...))
}

type Response interface {
	Write(io.Writer) error
	StatusCode() int
}

type responseBody struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data,omitempty"`
	Error  string          `json:"error,omitempty"`
}

type response struct {
	status Status
	data   interface{}
	err    error
}

func (res *response) Write(w io.Writer) error {
	var (
		b = &responseBody{
			Status: int(res.status),
		}
		err error
	)
	if res.status == StatusOK {
		b.Data, err = json.Marshal(res.data)
		if err != nil {
			return err
		}
	} else {
		b.Error = res.err.Error()
	}
	return json.NewEncoder(w).Encode(b)
}

func (res *response) StatusCode() int {
	return int(res.status)
}
