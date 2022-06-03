package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(h http.Handler, opts ...Option) *Server {
	hs := &http.Server{
		Handler: h,
	}

	s := &Server{
		server:          hs,
		notify:          make(chan error),
		shutdownTimeout: time.Second * 5,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	fmt.Println("httpserver start at", s.server.Addr)
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
