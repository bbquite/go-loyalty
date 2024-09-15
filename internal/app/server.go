package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *Config, mux *chi.Mux, logger *zap.SugaredLogger) error {

	s.httpServer = &http.Server{
		Addr:           cfg.RunAddress,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
