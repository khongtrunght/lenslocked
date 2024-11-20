package controllers

import (
	"fmt"
	"log/slog"
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
	user, err := u.UserService.Create(email, password)
	if err != nil {
		slog.Error("error creating user", slog.Any("error", err))
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Created user: %+v", user)
}
