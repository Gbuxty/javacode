package server

import (
	"context"
	"fmt"
	"java-code/config"
	"java-code/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
	cfg        config.HTTPServer
}

func NewServer(
	cfg config.HTTPServer,
	logger *logger.Logger,
	handler http.Handler,
) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s Server) ListenAndServe() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s Server) Run() {
	go func() {
		s.logger.Info("Starting HTTP server", zap.Int("port", s.cfg.Port))

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Listen:", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ServerShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Shutdown Server ...", zap.Error(err))
	}

	s.logger.Info("Server stopped.")
}
