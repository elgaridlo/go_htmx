package router

import (
	"net/http"
	"project/configs"
	"project/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func SetupRouter(
	cfg configs.Config,
	umw controllers.UserMiddleware,
	usersC controllers.Users,

) *chi.Mux {
	r := chi.NewRouter()

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	r.Use(csrfMw)
	r.Use(umw.SetUser)

	// Authentication
	r.Group(func(r chi.Router) {
		r.Get("/signin", usersC.SignIn)
		r.Post("/signin", usersC.ProcessSignIn)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.GetUsers)
		r.Get("/add", usersC.ViewCreateUser)
		r.Post("/", usersC.ProcessCreateUser)
		r.Get("/delete/{id}", usersC.DeleteUser)
		r.Get("/edit/{id}", usersC.ViewUpdateUser)
		r.Post("/edit", usersC.ProcessUpdateUser)
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	return r
}
