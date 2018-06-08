package controllers

import (
	"../views"
	"fmt"
	"net/http"
)

type Users struct {
	NewView *views.View
}

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	u.NewView.Render(w, nil)
	return
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:""`
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := decodeFrom(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}
