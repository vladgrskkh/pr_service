package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	srv    *http.Server
	logger *slog.Logger
}

func NewServer(srv *http.Server, logger *slog.Logger) *Server {
	return &Server{
		srv:    srv,
		logger: logger,
	}
}

func (s *Server) Serve() error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		// catch signals
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		sig := <-quit

		s.logger.Info("shutting down server", slog.String("signal", sig.String()))

		// for now timeout is set to 5 seconds (change for production)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return fmt.Errorf("something went wrong while shutting down: %w", err)
	}

	s.logger.Info("server stopped", slog.String("addr", s.srv.Addr))
	return nil
}
