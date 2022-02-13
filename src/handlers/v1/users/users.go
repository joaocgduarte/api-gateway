package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/plagioriginal/api-gateway/domain"
	"github.com/plagioriginal/api-gateway/helpers"
)

type UsersHandler struct {
	Logger        *log.Logger
	Validator     *validator.Validate
	UsersClient   domain.UsersClient
	CookieHandler domain.CookieHandler
}

func New(
	usersClient domain.UsersClient,
	cookieHandler domain.CookieHandler,
	v *validator.Validate,
	l *log.Logger,
) domain.UsersHttpHandler {
	return UsersHandler{
		UsersClient:   usersClient,
		CookieHandler: cookieHandler,
		Logger:        l,
		Validator:     v,
	}
}

func (uh UsersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	refreshToken := uh.CookieHandler.GetRefreshToken(r)

	if len(refreshToken) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		helpers.JSON(w, r, "invalid token")
		return
	}

	response, err := uh.UsersClient.Logout(ctx, refreshToken)

	if err != nil {
		uh.Logger.Printf("error on logout %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		helpers.JSON(w, r, "internal error")
		return
	}

	uh.CookieHandler.GenerateCookiesFromTokens(w, "", "")

	w.WriteHeader(http.StatusOK)
	helpers.JSON(w, r, response)
}

func (uh UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	request := domain.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		uh.Logger.Printf("login request body error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		helpers.JSON(w, r, "invalid request")
		return
	}

	err = uh.Validator.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors).Error()
		w.WriteHeader(http.StatusBadRequest)
		helpers.JSON(w, r, validationErrors)
		return
	}

	result, err := uh.UsersClient.Login(ctx, request)
	if err != nil {
		uh.Logger.Printf("error on client upon login: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		helpers.JSON(w, r, "internal error")
		return
	}

	uh.CookieHandler.GenerateCookiesFromTokens(w, result.AccessToken, result.RefreshToken)

	result.AccessToken = ""
	result.RefreshToken = ""

	w.WriteHeader(http.StatusOK)
	helpers.JSON(w, r, result)

}

func (uh UsersHandler) RefreshJWT(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	refreshToken := uh.CookieHandler.GetRefreshToken(r)
	if len(refreshToken) == 0 {
		w.WriteHeader(http.StatusNotFound)
		helpers.JSON(w, r, "invalid token")
		return
	}

	request := domain.RefreshRequest{
		RefreshToken: refreshToken,
	}

	err := uh.Validator.Var(request.RefreshToken, "required")

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		helpers.JSON(w, r, validationErrors)
		return
	}

	result, err := uh.UsersClient.RefreshJWT(ctx, request.RefreshToken)
	if err != nil {
		uh.Logger.Printf("error on client upon refresh: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		helpers.JSON(w, r, "error still don't check what it is (code)")
		return
	}

	uh.CookieHandler.GenerateCookiesFromTokens(w, result.AccessToken, result.RefreshToken)

	result.AccessToken = ""
	result.RefreshToken = ""

	w.WriteHeader(http.StatusOK)
	helpers.JSON(w, r, result)
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

	helpers.JSON(w, r, response)
}
