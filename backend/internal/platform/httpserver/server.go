package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func New(address string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              address,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		shutdownTimeout: 10 * time.Second,
	}
}

func (server *Server) Run(ctx context.Context) error {
	serverError := make(chan error, 1)

	go func() {
		log.Printf(
			"Starting SamaTalk API on http://localhost%s",
			server.httpServer.Addr,
		)

		serverError <- server.httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf(
			"listen and serve: %w",
			err,
		)

	case <-ctx.Done():
		log.Println(
			"shutdown signal received",
		)

		shutdownContext, cancelShutdown :=
			context.WithTimeout(
				context.Background(),
				server.shutdownTimeout,
			)

		defer cancelShutdown()

		log.Println(
			"shutting down HTTP server...",
		)

		err := server.httpServer.Shutdown(
			shutdownContext,
		)
		if err != nil {
			return fmt.Errorf(
				"shutdown HTTP server: %w",
				err,
			)
		}

		err = <-serverError

		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf(
				"HTTP server stopped unexpectedly: %w",
				err,
			)
		}

		log.Println(
			"SamaTalk API stopped successfully",
		)

		return nil
	}
}
