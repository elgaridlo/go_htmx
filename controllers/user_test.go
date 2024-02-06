package controllers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"project/controllers"
	"project/models"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupUserController(db *sql.DB) controllers.Users {
	// Set up services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	userMiddleware := &controllers.UserMiddleware{
		SessionService: sessionService,
	}

	// Set up controllers
	usersC := controllers.Users{
		UserService:       userService,
		SessionService:    sessionService,
		MiddlewareService: userMiddleware,
	}

	return usersC
}

func TestRenderViewSignIn(t *testing.T) {
	r := chi.NewRouter()
	usersC := setupUserController(db)
	r.Get("/signin", usersC.SignIn)

	req, err := http.NewRequest("GET", "/signin", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := execRequest(req, r)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

func TestProcesSignIn(t *testing.T) {
	r := chi.NewRouter()
	usersC := setupUserController(db)
	r.Post("/signin", usersC.ProcessSignIn)

	formData := url.Values{}
	formData.Set("email", "admin@gmail.com")
	formData.Set("password", "password123")

	reqBody := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "/signin", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := execRequest(req, r)
	assert.Equal(t, http.StatusFound, w.Code, "Expected status code 302")
}

func TestGetUser(t *testing.T) {
	r := chi.NewRouter()
	usersC1 := setupUserController(db)
	r.Use(usersC1.MiddlewareService.SetUser)
	r.Get("/users", usersC1.GetUsers)

	currentUserReq, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	userRes := execRequest(currentUserReq, r)
	assert.Equal(t, http.StatusOK, userRes.Code, "Expect status 200")
}

func TestUpdateUser(t *testing.T) {
	r := chi.NewRouter()
	usersC1 := setupUserController(db)

	r.Post("/signin", usersC1.ProcessSignIn)
	r.Get("/users/edit/{id}", usersC1.ViewUpdateUser)
	r.With(usersC1.MiddlewareService.SetUser, usersC1.MiddlewareService.RequireUser).Post("/users/edit", usersC1.ProcessUpdateUser)

	formData := url.Values{}
	formData.Set("email", "staff@gmail.com")
	formData.Set("password", "password123")

	reqBody := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "/signin", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := execRequest(req, r)
	cookie := w.Result().Cookies()
	assert.Equal(t, http.StatusFound, w.Code, "Expected status code 302")

	req, err = http.NewRequest("GET", "/users/edit/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	w = execRequest(req, r)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	formEdit := url.Values{}
	formEdit.Set("id", "3")
	formEdit.Set("email", "staff@gmail.com")
	formEdit.Set("name", "staff khusus")

	editBody := strings.NewReader(formEdit.Encode())

	currentUserReq, err := http.NewRequest("POST", "/users/edit", editBody)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cookie {
		currentUserReq.AddCookie(c)
	}

	currentUserReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	userRes := execRequest(currentUserReq, r)
	assert.Equal(t, http.StatusFound, userRes.Code, "Expected status 302")

	formEdit = url.Values{}
	formEdit.Set("id", "3")
	formEdit.Set("email", "staff@gmail.com")
	formEdit.Set("name", "staff")

	editBody = strings.NewReader(formEdit.Encode())

	currentUserReq, err = http.NewRequest("POST", "/users/edit", editBody)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cookie {
		currentUserReq.AddCookie(c)
	}

	currentUserReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	userRes = execRequest(currentUserReq, r)
	assert.Equal(t, http.StatusFound, userRes.Code, "Expected status 302")
}

func TestDeleteUser(t *testing.T) {
	r := chi.NewRouter()
	controller := setupUserController(db)

	newUser := models.User{
		Email:        "Sudarmono",
		PasswordHash: "passwordkita",
	}

	user, err := controller.UserService.CreateUser(newUser)
	if err != nil {
		t.Fatal(err)
	}

	r.Delete("/user/{id}", controller.DeleteUser)

	deleteReq, err := http.NewRequest("DELETE", "/user/"+strconv.Itoa(user.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	deleteReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	deleteRes := execRequest(deleteReq, r)
	redirect, _ := deleteRes.Result().Location()

	assert.Equal(t, "/users", redirect.Path, "Expected redirect to /users")
	assert.Equal(t, http.StatusFound, deleteRes.Code, "Expected status code 404")
}

func TestCreateUser(t *testing.T) {
	r := chi.NewRouter()

	userC := setupUserController(db)

	cookie := loginUser(t, r, "admin@gmail.com", "password123")
	r.With(userC.MiddlewareService.SetUser,
		userC.MiddlewareService.RequireUser).Post("/users", userC.ProcessCreateUser)

	formPost := url.Values{}
	formPost.Set("email", "markus_horizon@gmail.com")
	formPost.Set("password", "password123")
	formPost.Set("name", "markus horizontal")

	postBody := strings.NewReader(formPost.Encode())

	createReq, err := http.NewRequest("POST", "/users", postBody)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cookie {
		createReq.AddCookie(c)
	}

	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createRes := execRequest(createReq, r)
	assert.Equal(t, http.StatusCreated, createRes.Code, "Expected status 201")

	cookie = loginUser(t, r, "markus_horizon@gmail.com", "password123")

	assert.NotEmpty(t, cookie, "Expected not empty cookie")

	var id int

	row, err := userC.UserService.DB.Query(`
	SELECT id FROM users where email = $1 limit 1
	`, "markus_horizon@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	if row.Next() {
		err = row.Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No rows found")
	}

	r.Delete("/user/{id}", userC.DeleteUser)

	deleteReq, err := http.NewRequest("DELETE", "/user/"+strconv.Itoa(id), nil)
	if err != nil {
		t.Fatal(err)
	}
	deleteReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	deleteRes := execRequest(deleteReq, r)
	redirect, _ := deleteRes.Result().Location()

	assert.Equal(t, "/users", redirect.Path, "Expected redirect to /users")
	assert.Equal(t, http.StatusFound, deleteRes.Code, "Expected status code 404")
}

func execRequest(req *http.Request, r *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}
