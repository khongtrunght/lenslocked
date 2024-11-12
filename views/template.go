package views

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}
	return Template{HtmlTpl: tpl}, nil
}

type Template struct {
	HtmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.HtmlTpl.Execute(w, data)
	if err != nil {
		slog.Error("error executing template", slog.Any("error", err))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
