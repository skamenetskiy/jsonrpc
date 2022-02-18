package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewClient(addr string) Client {
	return &client{addr}
}

type Client interface {
	Call(method string, requestData, responseData interface{}) error
}

type client struct {
	addr string
}

func (c *client) Call(method string, requestData, responseData interface{}) error {
	var (
		req = &requestBody{
			Method: method,
		}
		err error
	)
	req.Data, err = json.Marshal(requestData)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if err = json.NewEncoder(buf).Encode(req); err != nil {
		return err
	}
	resp, err := http.Post(c.addr, contentType, buf)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	res := new(responseBody)
	if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
		return err
	}
	if res.Status != int(StatusOK) {
		return fmt.Errorf(res.Error)
	}
	return json.Unmarshal(res.Data, responseData)
}
