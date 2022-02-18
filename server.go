package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func NewServer(name string) Server {
	return &server{
		name:     name,
		handlers: make(map[string]Handler, 0),
	}
}

type Server interface {
	Handle(method string, handler Handler)
	Listen(addr string) error
	Shutdown()
}

type Handler func(req Request) Response

type server struct {
	name     string
	handlers map[string]Handler
	mu       sync.Mutex
	http     *http.Server
}

func (srv *server) Handle(method string, handler Handler) {
	srv.mu.Lock()
	if _, exists := srv.handlers[method]; exists {
		log.Panicf("handler %s already exists", method)
	}
	srv.handlers[method] = handler
	srv.mu.Unlock()
}

func (srv *server) Listen(addr string) error {
	srv.http = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(srv.handle),
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-shutdown
		log.Printf("received %s, shutting down\n", sig.String())

	}()
	return srv.http.ListenAndServe()
}

func (srv *server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.http.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown gracefully: %s", err.Error())
	}
}

func (srv *server) handle(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	reqBody := new(requestBody)
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		srv.write(w, Error(fmt.Errorf("failed to parse request: %s", err.Error())))
		return
	}
	var (
		req = &request{
			req:  r,
			body: reqBody,
		}
		res Response
	)
	if handler, ok := srv.handlers[req.body.Method]; ok {
		res = handler(req)
	} else {
		res = Errorf("method %s not found", req.body.Method)
	}
	srv.write(w, res)
}

func (srv *server) write(w http.ResponseWriter, res Response) {
	w.Header().Set("content-type", contentType)
	w.Header().Set("server", srv.name)
	w.WriteHeader(res.StatusCode())
	if err := res.Write(w); err != nil {
		log.Printf("failed to write: %s\n", err.Error())
	}
}
