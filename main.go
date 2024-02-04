package main

import (
	"fmt"
	"net/http"
	configs "project/config"
	"project/controllers"
	"project/migrations"
	"project/models"
	templates "project/template"
	views "project/view"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := configs.LoadEnvConfig(".env")
	if err != nil {
		panic(err)
	}

	// Setup the database
	db := configs.SetupDatabase(cfg.PSQL)
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Set up services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}

	// Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: &models.SessionService{},
	}
	r := chi.NewRouter()

	// Set up controllers
	usersC := controllers.Users{
		UserService:    userService,
		SessionService: sessionService,
	}
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "sign.html", "index.html"))

	// Dashboard
	r.With(umw.RequireUser).Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"sign.html", "index.html",
	))))

	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
