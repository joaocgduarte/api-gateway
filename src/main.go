package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/plagioriginal/api-gateway/handlers"
	"github.com/plagioriginal/api-gateway/middlewares"
)

// Inits the Http Server
func initHttpServer(r http.Handler, logger *log.Logger) {
	appPort := os.Getenv("API_PORT")

	server := &http.Server{
		Addr:    ":" + appPort,
		Handler: r,
	}

	// Go routine to begin the server
	go func() {
		logger.Printf("Listening to port %s\n", server.Addr)
		err := server.ListenAndServe()

		if err != nil {
			logger.Fatalln(err)
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

	log.Println("Shutting down server...")

	if err := server.Shutdown(timeoutContext); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func main() {
	logger := log.New(os.Stdout, "api-gateway: ", log.Flags())

	r := chi.NewRouter()
	userHandler := handlers.NewUserHandler(logger, nil)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger, NoColor: false}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares.SetJsonContentType)

	r.Route("/", func(r chi.Router) {
		r.Get("/", userHandler.Hello)
	})

	initHttpServer(r, logger)
}
