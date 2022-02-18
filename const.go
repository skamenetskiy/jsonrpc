package jsonrpc

import "net/http"

type Status int

const (
	contentType = "application/json"

	StatusOK    Status = http.StatusOK
	StatusError Status = http.StatusInternalServerError
)
