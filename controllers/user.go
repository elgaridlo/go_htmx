package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"project/context"
	"project/models"
	templates "project/template"
	"project/utils"
	views "project/view"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Users struct {
	Templates struct {
		SignIn Template
	}
	UserService       *models.UserService
	SessionService    *models.SessionService
	MiddlewareService *UserMiddleware
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	views.Must(
		views.ParseFS(
			templates.FS,
			"index.html",
			"sign.html",
		),
	).Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users", http.StatusFound)
}

func (u Users) GetUsers(w http.ResponseWriter, r *http.Request) {
	paginationQuery := utils.PaginationParams{}
	queryName := r.URL.Query().Get("name")
	paginationQuery.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	paginationQuery.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	users, err := u.UserService.UserMany(queryName, paginationQuery)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	views.Must(
		views.ParseFS(
			templates.FS,
			"index.html",
			"pages/users.html",
		),
	).Execute(w, r, users)
}

func (u Users) ViewCreateUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Name     string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Name = r.FormValue("name")
	data.Password = r.FormValue("password")

	views.Must(
		views.ParseFS(
			templates.FS,
			"index.html",
			"pages/add.html",
		),
	).Execute(w, r, data)
}

func (u Users) ProcessCreateUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Name     string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Name = r.FormValue("name")
	data.Password = r.FormValue("password")

	postUser := models.User{
		Email:        data.Email,
		Name:         data.Name,
		PasswordHash: data.Password,
	}

	_, err := u.UserService.CreateUser(postUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusCreated)
}

func (u Users) ViewUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	users, err := u.UserService.UserDetail(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return
	}

	views.Must(
		views.ParseFS(
			templates.FS,
			"index.html",
			"pages/edit.html",
		),
	).Execute(w, r, users)
}

func (u Users) ProcessUpdateUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieSession)
	if err != nil {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}
	sessionToken := cookie.Value

	var data struct {
		Id    int
		Name  string
		Email string
	}

	data.Id, _ = strconv.Atoi(r.FormValue("id"))
	data.Name = r.FormValue("name")
	data.Email = r.FormValue("email")

	postUser := models.User{
		ID:    data.Id,
		Name:  data.Name,
		Email: data.Email,
	}

	_, err = u.UserService.UserUpdate(postUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, sessionToken)
	http.Redirect(w, r, "/users", http.StatusFound)
}

func (u Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	err = u.UserService.UserDelete(idInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}
