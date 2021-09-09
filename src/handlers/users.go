package handlers

import (
	"log"
	"net/http"

	"github.com/plagioriginal/api-gateway/domain"
)

type UsersHandler struct {
	BaseHandler
	l          *log.Logger
	ApiService *domain.APIService
}

func NewUserHandler(l *log.Logger, apiService *domain.APIService) UsersHandler {
	return UsersHandler{l: l, ApiService: apiService}
}

func (uh UsersHandler) Hello(w http.ResponseWriter, r *http.Request) {
	uh.l.Println("Accessing route hello.")

	response := map[string]string{
		"hello": "world",
	}

	uh.JSON(w, r, response)
}
