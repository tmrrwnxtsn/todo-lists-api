package apiserver

import (
	"context"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/config"
	"net/http"
	"time"
)

type Server struct {
	config     *config.Config
	handler    http.Handler
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		config:  cfg,
		handler: handler,
	}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           s.config.BindAddr,
		Handler:        s.handler,
		MaxHeaderBytes: s.config.MaxHeaderBytes << 20, // MB
		ReadTimeout:    time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.config.WriteTimeout) * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
