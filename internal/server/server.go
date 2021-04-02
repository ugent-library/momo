package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	handler http.Handler
	host    string
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

func WithHost(h string) option {
	return func(s *Server) {
		s.host = h
	}
}

func WithPort(p int) option {
	return func(s *Server) {
		s.port = p
	}
}

func (s *Server) Start() {
	// mostly taken from examples here:
	// https://marcofranssen.nl/go-webserver-with-graceful-shutdown/
	// and here:
	// https://github.com/gorilla/mux#graceful-shutdown

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.handler,
		// set timeouts to avoid slowloris attacks.
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      1 * time.Minute,
	}

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	signalC := make(chan os.Signal, 1)

	signal.Notify(
		signalC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	log.Printf("Server is listening at http://%s\n", addr)

	// block until signal received
	<-signalC

	log.Println("Stopping...")

	// terminate with os.Exit(1) after second signal
	go func() {
		<-signalC
		log.Fatal("Kill - terminating...\n")
	}()

	// create a deadline
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()

	// stop immediately if there are no connections or
	// wait until deadline for connections to finish
	if err := srv.Shutdown(stopCtx); err != nil {
		log.Fatalf("Error stopping server: %v\n", err)
	} else {
		log.Println("Gracefully stopped")
	}

	os.Exit(0)
}
