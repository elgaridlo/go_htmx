package router

import (
	// "fmt"
	// "io/fs"
	// "net/http"
	// "os"
	configs "project/config"
	// "project/internal/controllers"
	// "project/public"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func SetupRouter(
	cfg configs.Config,

) *chi.Mux {
	r := chi.NewRouter()

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	r.Use(csrfMw)

	return r
}
