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
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/securecookie"
	usersClient "github.com/plagioriginal/api-gateway/clients/users"
	"github.com/plagioriginal/api-gateway/cookies"
	"github.com/plagioriginal/api-gateway/domain"
	usersHandler "github.com/plagioriginal/api-gateway/handlers/v1/users"
	"github.com/plagioriginal/api-gateway/middlewares"
	v1 "github.com/plagioriginal/api-gateway/router/v1"
	"github.com/plagioriginal/api-gateway/tokens"
	usersGrpc "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := log.New(os.Stdout, "api-gateway: ", log.Flags())
	r := chi.NewRouter()

	conn, err := grpc.Dial(os.Getenv("USERS_SERVICE_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalln(err)
	}
	defer conn.Close()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger, NoColor: false}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares.SetJsonContentType)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Set-Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	cookieEncoder := generateCookieHandler(logger)
	validator := validator.New()
	timeoutContext := time.Duration(3) * time.Second
	userClient := usersClient.New(usersGrpc.NewUsersClient(conn), logger, timeoutContext)
	tokenManager := tokens.NewTokenManager(os.Getenv("JWT_GENERATOR_SECRET"))
	usersHandler := usersHandler.New(userClient, cookieEncoder, validator, logger)
	authMiddleware := middlewares.NewAuthorizationMiddleware(tokenManager, userClient, cookieEncoder, logger)

	v1.New(
		"/v1",
		usersHandler,
		authMiddleware.RequireToken([]string{"admin"}),
	).GenerateRoutes(r)

	server := &http.Server{
		Addr:    ":" + os.Getenv("API_PORT"),
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
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContext)
	defer cancel()

	logger.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func generateCookieHandler(l *log.Logger) domain.CookieHandler {
	hashKey := securecookie.GenerateRandomKey(32)
	if hashKey == nil {
		l.Fatalln("couldn't generate hashkey for cookies")
	}
	blockKey := securecookie.GenerateRandomKey(24)
	if blockKey == nil {
		l.Fatalln("couldn't generate blockkey for cookies")
	}

	cookieEncoder := securecookie.New(hashKey, blockKey)
	return cookies.New(cookieEncoder)
}
