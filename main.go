package main

import (
	"fmt"
	"net/http"
	"project/cmd/migrations"
	"project/configs"
	"project/controllers"
	"project/models"
	"project/router"
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
		SessionService: sessionService,
	}

	// Set up controllers
	usersC := controllers.Users{
		UserService:    userService,
		SessionService: sessionService,
	}

	r := router.SetupRouter(
		cfg,
		umw,
		usersC,
	)

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
