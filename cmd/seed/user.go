package main

import (
	"fmt"
	"project/models"
)

func SeedUsers(services Services) {
	fmt.Println("Start seeding users...")

	var totalCount int

	row := services.UserService.DB.QueryRow("SELECT COUNT(*) as total FROM users")
	err := row.Scan(&totalCount)
	if err != nil {
		fmt.Println("Error scanning row:", err)
		return
	}

	if totalCount > 0 {
		fmt.Println("Seeding users already exist...")
		return
	}

	insUsers := []models.User{
		{Email: "admin@gmail.com", PasswordHash: "password123"},
		{Email: "manager@gmail.com", PasswordHash: "password123"},
		{Email: "staff@gmail.com", PasswordHash: "password123"},
	}

	for _, user := range insUsers {
		_, err := services.UserService.Create(user)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Seeding users completed...")
}
