package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plagioriginal/api-gateway/domain"
)

type Router struct {
	prefix              string
	usersHandler        domain.UsersHttpHandler
	adminAuthMiddleware func(next http.Handler) http.Handler
}

func New(
	prefix string,
	usersHandler domain.UsersHttpHandler,
	adminAuthMiddleware func(next http.Handler) http.Handler,
) Router {
	return Router{
		prefix:              prefix,
		usersHandler:        usersHandler,
		adminAuthMiddleware: adminAuthMiddleware,
	}
}

func (router Router) GenerateRoutes(mux *chi.Mux) {
	mux.Route(router.prefix+"/users", func(r chi.Router) {
		r.Post("/login", router.usersHandler.Login)
		r.Post("/refresh", router.usersHandler.RefreshJWT)
		r.Post("/logout", router.usersHandler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(router.adminAuthMiddleware)
			r.Get("/", router.usersHandler.AddUser)
		})
	})
}
