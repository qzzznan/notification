package httpserver

import (
	"context"
	"net"
)

type Option func(*Server)

func WithContext(ctx context.Context) Option {
	return func(s *Server) {
		s.server.BaseContext = func(listener net.Listener) context.Context {
			return ctx
		}
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.server.Addr = addr
	}
}
