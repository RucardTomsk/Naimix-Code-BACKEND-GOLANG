package server

import (
	"context"
	"fmt"
	"net/http"
)

// Server is a custom http server.
type Server struct {
	httpServer *http.Server
}

// Run starts everlasting listener on provided host and port.
func (s *Server) Run(host string, port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown shuts off running server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
