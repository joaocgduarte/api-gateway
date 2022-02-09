package main

import (
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/plagioriginal/api-gateway/api/service"
	"github.com/plagioriginal/api-gateway/handlers"
	"github.com/plagioriginal/api-gateway/middlewares"
	"github.com/plagioriginal/api-gateway/server"
	"github.com/plagioriginal/api-gateway/tokens"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc"
)

func main() {
	logger := log.New(os.Stdout, "api-gateway: ", log.Flags())
	r := chi.NewRouter()

	userServiceHost := os.Getenv("USERS_SERVICE_HOST")
	if len(userServiceHost) == 0 {
		userServiceHost = "users-service:8080"
	}
	conn, err := grpc.Dial(userServiceHost, grpc.WithInsecure())
	if err != nil {
		logger.Fatalln(err)
	}
	defer conn.Close()

	userClient := users.NewUsersClient(conn)
	validator := validator.New()

	timeoutContext := time.Duration(2) * time.Second
	apiService := service.New(userClient, logger, timeoutContext)
	tokenManager := tokens.NewTokenManager(os.Getenv("JWT_GENERATOR_SECRET"))

	authMiddleware := middlewares.NewAuthorizationMiddleware(tokenManager, userClient, logger)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger, NoColor: false}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares.SetJsonContentType)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Set-Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	handlers.NewUserHandler(apiService, validator, logger, authMiddleware, r)

	serverInitializer := server.ServerInitializer{
		Logger:  logger,
		Handler: r,
		Port:    os.Getenv("API_PORT"),
	}

	serverInitializer.Init()
}
