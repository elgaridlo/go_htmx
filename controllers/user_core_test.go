package controllers_test

import (
	"database/sql"
	"project/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupCoreUser(db *sql.DB) *models.UserService {
	userService := &models.UserService{
		DB: db,
	}

	return userService
}

func TestUnitUserCRUD(t *testing.T) {
	userService := setupCoreUser(db)
	name := ""

	paginationParams := struct {
		Page     int
		PageSize int
	}{
		Page:     1,
		PageSize: 1,
	}

	users, err := userService.UserMany(name, paginationParams)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(users.Data), 1)
	for _, div := range users.Data {
		assert.Equal(t, div.ID, 1)
	}

	insertUser := models.User{
		Name:         "Sudarmaji",
		Email:        "sudarmaji@gmail.com",
		PasswordHash: "password321",
	}

	create, err := userService.CreateUser(insertUser)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, create.ID, "Expected to have ID")

	expectedRes := models.User{
		ID:    create.ID,
		Name:  insertUser.Name,
		Email: strings.ToLower(insertUser.Email),
	}

	assert.Equal(t, expectedRes.Email, create.Email, "Expected to be same value")
	assert.Equal(t, expectedRes.Name, create.Name, "Expected to be same value")

	updateUser := models.User{
		ID:    create.ID,
		Name:  "Sudarmaji Pramono",
		Email: "sudarmaji@gmail.com",
	}

	update, err := userService.UserUpdate(updateUser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &updateUser, update, "Expected update to be same value")

	deleteUser := userService.UserDelete(create.ID)
	assert.Nil(t, deleteUser, "Expected to be nil")
}

func TestAuthenticate(t *testing.T) {
	userService := setupCoreUser(db)

	auth, err := userService.Authenticate("admin@gmail.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, auth.ID, "Expected not empty")

	_, err = userService.Authenticate("admin1@gmail.com", "password123")

	assert.Contains(t, err.Error(), "authenticate: sql: no rows in result set", "Expected no rows in result set")
}
