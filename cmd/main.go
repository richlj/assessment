package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/richlj/sliide/go-challenge-master/pkg"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "the TCP address for the server to listen on, in the form 'host:port'")

	// pkg.App

	// app gets initialised with configuration.
	// as an example we've added 3 providers and a default configuration

	app = pkg.App{
		ContentClients: map[pkg.Provider]pkg.Client{
			pkg.Provider1: pkg.SampleContentProvider{Source: pkg.Provider1},
			pkg.Provider2: pkg.SampleContentProvider{Source: pkg.Provider2},
			pkg.Provider3: pkg.SampleContentProvider{Source: pkg.Provider3},
		},
		Config: pkg.DefaultConfig,
	}
)

func main() {
	log.Printf("initalising server on %s", *addr)

	srv := http.Server{
		Addr:    *addr,
		Handler: app,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
