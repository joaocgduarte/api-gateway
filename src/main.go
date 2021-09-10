package main

import (
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/plagioriginal/api-gateway/api/service"
	"github.com/plagioriginal/api-gateway/handlers"
	"github.com/plagioriginal/api-gateway/middlewares"
	"github.com/plagioriginal/api-gateway/protos/protos"
	"github.com/plagioriginal/api-gateway/server"
	"google.golang.org/grpc"
)

func main() {
	logger := log.New(os.Stdout, "api-gateway: ", log.Flags())
	r := chi.NewRouter()

	conn, err := grpc.Dial("api-users-todos:8080", grpc.WithInsecure())
	if err != nil {
		logger.Fatalln(err)
	}
	defer conn.Close()

	userClient := protos.NewUsersClient(conn)
	validator := validator.New()

	timeoutContext := time.Duration(2) * time.Second
	apiService := service.New(userClient, logger, timeoutContext)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger, NoColor: false}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares.SetJsonContentType)

	handlers.NewUserHandler(apiService, validator, logger, r)

	serverInitializer := server.ServerInitializer{
		Logger:  logger,
		Handler: r,
		Port:    os.Getenv("API_PORT"),
	}

	serverInitializer.Init()
}
