package controllers_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"
	"project/configs"
	"project/models"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var (
	db  *sql.DB
	cfg configs.Config
)

func TestMain(m *testing.M) {
	var err error
	cfg, err = configs.LoadEnvConfig("../.env")
	if err != nil {
		log.Fatal(err)
	}

	db, err = models.Open(cfg.PSQL)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func loginUser(t *testing.T, router *chi.Mux, email, password string) []*http.Cookie {
	usersC := setupUserController(db)

	router.Post("/signin", usersC.ProcessSignIn)
	formData := url.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	reqBody := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "/signin", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := execRequest(req, router)
	cookies := w.Result().Cookies()

	assert.Equal(t, http.StatusFound, w.Code, "Expected status code 302")

	return cookies
}
