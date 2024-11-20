package controllers

import (
	"fmt"
	"net/http"

	"github.com/khongtrunght/lenslocked/models"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil, nil)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Fprintf(w, "Email: %s\n", email)
	fmt.Fprintf(w, "Password: %s\n", password)
}
