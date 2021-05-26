package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	handler http.Handler
	port    int
}

type option func(*Server)

func New(h http.Handler, opts ...option) *Server {
	s := &Server{handler: h}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithPort(p int) option {
	return func(s *Server) {
		s.port = p
	}
}

func (s *Server) Start() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.handler,
		// Set timeouts to avoid slowloris attacks.
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      1 * time.Minute,
	}

	// Listen on a different Goroutine so the application doesn't stop here.
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	log.Printf("Server is listening at %s\n", srv.Addr)

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println("Shutting down gracefully, press Ctrl+C again to force")

	// Shutdown with a maximum timeout of 10 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
