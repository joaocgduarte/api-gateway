package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Constructor to initialize the server.
type ServerInitializer struct {
	Handler http.Handler
	Logger  *log.Logger
	Port    string
}

// Inits the http server.
func (si ServerInitializer) Init() {
	server := &http.Server{
		Addr:    ":" + si.Port,
		Handler: si.Handler,
	}

	// Go routine to begin the server
	go func() {
		si.Logger.Printf("Listening to port %s\n", server.Addr)
		err := server.ListenAndServe()

		if err != nil {
			si.Logger.Fatalln(err)
		}
	}()

	// Wait for an interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Attempt a graceful shutdown
	timeoutContext, cancel := context.
		WithTimeout(
			context.Background(),
			time.Duration(2)*time.Second,
		)
	defer cancel()

	si.Logger.Println("Shutting down server...")

	if err := server.Shutdown(timeoutContext); err != nil {
		si.Logger.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
