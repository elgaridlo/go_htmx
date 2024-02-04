package main

import (
	"database/sql"
	"fmt"
	"project/configs"
	"project/models"
)

type Services struct {
	PSQL        models.PostgresConfig
	UserService *models.UserService
}

func main() {
	_, err := configs.LoadEnvConfig(".env")
	if err != nil {
		panic(err)
	}

	var services Services

	services.PSQL = models.DefaultPostgresConfig()
	fmt.Println("Connecting to db at", services.PSQL.Host)

	db, err := models.Open(services.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	initServices(&services, db)

	fmt.Println("Starting seeders application.")
	SeedUsers(services)
	fmt.Println("Process seeding completed.")
}

func initServices(services *Services, db *sql.DB) {
	services.UserService = &models.UserService{
		DB: db,
	}
}
