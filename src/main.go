package main

import (
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/plagioriginal/api-gateway/handlers"
	"github.com/plagioriginal/api-gateway/middlewares"
	"github.com/plagioriginal/api-gateway/server"
)

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

	serverInitializer := server.ServerInitializer{
		Logger:  logger,
		Handler: r,
		Port:    os.Getenv("API_PORT"),
	}

	serverInitializer.Init()
}
