package controllers

import (
	"net/http"

	"github.com/a-h/templ"
)

func StaticHandler(tpl templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
