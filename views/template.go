package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/khongtrunght/lenslocked/context"
	"github.com/khongtrunght/lenslocked/models"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(patterns[0]))
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (string, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
		"errors": func() []string {
			return []string{}
		},
	})

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}
	return Template{HtmlTpl: tpl}, nil
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.HtmlTpl.Clone()
	if err != nil {
		slog.Error("error cloning template", slog.Any("error", err))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			return errMsgs
		},
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		slog.Error("error executing template", slog.Any("error", err))
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		if err != nil {
			var publicErr public
			if errors.As(err, &publicErr) {
				msgs = append(msgs, publicErr.Public())
			} else {
				msgs = append(msgs, "Something went wrong.")
			}
		}
	}
	return msgs
}
