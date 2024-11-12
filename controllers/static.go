package controllers

import (
	"html/template"
	"net/http"

	"github.com/khongtrunght/lenslocked/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days.",
		},
		{
			Question: "Can I cancel anytime?",
			Answer:   "Yes! You can cancel anytime without any fees.",
		},
		{
			Question: "Is there a mobile app?",
			Answer:   "Yes! We have a mobile app for both iOS and Android.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us at <a href="mailto:khongtrunght@gmail.com">khongturnght@gmail.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil, questions)
	}
}
