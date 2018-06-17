package controllers

import (
	"../models"
	"../views"
	"fmt"
	"net/http"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	US        *models.UserService
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		US:        us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	u.NewView.Render(w, nil)
	return
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := decodeFrom(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Email:    form.Email,
		Name:     form.Name,
		Password: form.Password,
	}

	if err := u.US.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	signIn(w, &user)
	http.Redirect(w, r, "/cookie_test", http.StatusFound)
	return
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := decodeFrom(r, &form); err != nil {
		panic(err)
	}

	user, err := u.US.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrorInvalidPassword:
		fmt.Fprintln(w, "Invalid Password Provided.")
		return
	case models.ErrorNotFound:
		fmt.Fprintln(w, "Invalid email address.")
		return
	case nil:
		signIn(w, user)
		http.Redirect(w, r, "/cookie_test", http.StatusFound)
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Email is: ", cookie.Value)
	fmt.Fprintln(w, cookie)

	return
}

func signIn(w http.ResponseWriter, u *models.User) {
	cookie := http.Cookie{
		Name:  "email",
		Value: u.Email,
	}
	http.SetCookie(w, &cookie)
}
