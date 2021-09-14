package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/plagioriginal/api-gateway/domain"
	"github.com/plagioriginal/api-gateway/http_renderer"
	"github.com/plagioriginal/api-gateway/middlewares"
)

type UsersHandler struct {
	Logger         *log.Logger
	Validator      *validator.Validate
	ApiService     domain.APIService
	AuthMiddleware middlewares.AuthorizationMiddleware
}

func NewUserHandler(apiService domain.APIService, v *validator.Validate, l *log.Logger, aw middlewares.AuthorizationMiddleware, r *chi.Mux) {
	handler := UsersHandler{ApiService: apiService, Logger: l, Validator: v, AuthMiddleware: aw}

	r.Route("/users", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.RefreshJWT)
		r.Post("/logout", handler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(aw.RequireAdminValidToken)
			r.Get("/", handler.AddUser)
		})
	})
}

func (uh UsersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	tokenCookie, err := r.Cookie("refresh-token")

	if err != nil {
		uh.Logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: err.Error()})
		return
	}

	refreshToken := tokenCookie.Value

	response, err := uh.ApiService.Logout(ctx, refreshToken)

	if err != nil {
		uh.Logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: err.Error()})
		return
	}

	accessCookies := []http.Cookie{
		{
			Name:     "access-token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			// SameSite: http.SameSiteNoneMode,
			Expires: time.Now().Add(time.Hour * 24 * 14),
			// Secure:   true,
		},
		{
			Name:     "refresh-token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			// SameSite: http.SameSiteNoneMode,
			Expires: time.Now().Add(time.Hour * 24 * 14),
			// Secure:   true,
		},
	}

	for _, cookie := range accessCookies {
		http.SetCookie(w, &cookie)
	}

	w.WriteHeader(http.StatusOK)
	http_renderer.JSON(w, r, response)
}

func (uh UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	request := domain.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		uh.Logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: err.Error()})
		return
	}

	err = uh.Validator.Struct(request)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors).Error()
		w.WriteHeader(http.StatusBadRequest)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: validationErrors})
		return
	}

	result, err := uh.ApiService.Login(ctx, request)

	if err != nil {
		uh.Logger.Println(err)
		w.WriteHeader(http.StatusNotFound)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: err.Error()})
		return
	}

	accessCookies := []http.Cookie{
		{
			Name:     "access-token",
			Value:    result.AccessToken,
			HttpOnly: true,
			Path:     "/",
			// SameSite: http.SameSiteNoneMode,
			Expires: time.Now().Add(time.Hour * 24 * 14),
			// Secure:   true,
		},
		{
			Name:     "refresh-token",
			Value:    result.RefreshToken,
			HttpOnly: true,
			Path:     "/",
			// SameSite: http.SameSiteNoneMode,
			Expires: time.Now().Add(time.Hour * 24 * 14),
			// Secure:   true,
		},
	}

	for _, cookie := range accessCookies {
		http.SetCookie(w, &cookie)
	}

	w.WriteHeader(http.StatusOK)
	http_renderer.JSON(w, r, result)

}

func (uh UsersHandler) RefreshJWT(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	refreshToken, err := r.Cookie("refresh-token")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: err.Error()})
		return
	}

	request := domain.RefreshRequest{
		RefreshToken: refreshToken.Value,
	}

	err = uh.Validator.Var(request.RefreshToken, "required")

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		uh.Logger.Panicln(validationErrors)
		w.WriteHeader(http.StatusBadRequest)
		http_renderer.JSON(w, r, validationErrors)
		return
	}

	result, err := uh.ApiService.RefreshJWT(ctx, request.RefreshToken)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http_renderer.JSON(w, r, domain.FailRestResponse{Errors: err.Error()})
		return
	}

	accessCookies := []http.Cookie{
		{
			Name:     "access-token",
			Value:    result.AccessToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 14),
		},
		{
			Name:     "refresh-token",
			Value:    result.RefreshToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 14),
		},
	}

	for _, cookie := range accessCookies {
		http.SetCookie(w, &cookie)
	}

	w.WriteHeader(http.StatusOK)
	http_renderer.JSON(w, r, result)
}

func (uh UsersHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := ctx.Value("userId")
	userRole := ctx.Value("userRole")

	uh.Logger.Println(userId)
	uh.Logger.Println(userRole)

	response := map[string]string{
		"hello": "world",
	}

	http_renderer.JSON(w, r, response)
}

func (uh UsersHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"hello": "world",
	}

	http_renderer.JSON(w, r, response)
}
