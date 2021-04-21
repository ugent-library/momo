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

	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	handler http.Handler
	host    string
	port    int
	ssl     bool
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

func WithSSL(b bool) option {
	return func(s *Server) {
		s.ssl = b
	}
}

func (s *Server) Start() {
	// graceful shutdown mostly taken from examples here:
	// https://marcofranssen.nl/go-webserver-with-graceful-shutdown/
	// https://github.com/gorilla/mux#graceful-shutdown
	// autocert mostly taken from examples here:
	// https://gist.github.com/samthor/5ff8cfac1f80b03dfe5a9be62b29d7f2
	// https://geekbrit.org/content/28388
	// https://blog.kowalczyk.info/article/Jl3G/https-for-free-in-go.html

	var srv *http.Server

	if s.ssl {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("letsencrypt"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.host),
		}
		srv = &http.Server{
			Addr:      ":https",
			TLSConfig: m.TLSConfig(),
			Handler:   s.handler,
			// set timeouts to avoid slowloris attacks.
			ReadHeaderTimeout: 20 * time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      1 * time.Minute,
		}

		// start http server
		go func() {
			// serve HTTP, which will redirect automatically to HTTPS
			h := m.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":http", h))
		}()

		// start https server
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				log.Fatalf("Error starting server: %v\n", err)
			}
		}()

		log.Printf("Server is listening at https://%s\n", s.host)
	} else {
		addr := fmt.Sprintf("%s:%d", s.host, s.port)

		srv = &http.Server{
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

		log.Printf("Server is listening at http://%s\n", addr)
	}

	signalC := make(chan os.Signal, 1)

	signal.Notify(
		signalC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

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
