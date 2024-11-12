package controllers

import (
	"net/http"

	"github.com/khongtrunght/lenslocked/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil, nil)
	}
}
