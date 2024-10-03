package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/khongtrunght/lenslocked/controllers"
	"github.com/khongtrunght/lenslocked/templates/layouts"
	"github.com/khongtrunght/lenslocked/templates/pages"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", controllers.StaticHandler(pages.Home()))
	router.Get("/about", controllers.StaticHandler(pages.About()))
	router.Get("/contact", controllers.StaticHandler(pages.Contact()))
	router.Get("/faq", controllers.StaticHandler(pages.Faq(
		[]pages.Question{
			{"What is LensLocked?", "LensLocked is a photo gallery website."},
			{"How do I contact you?", "You can contact us by visiting the contact page."},
			{"How do I see the photos?", "You can see the photos by visiting the gallery page."},
		},
	)))
	router.Get("/main", controllers.StaticHandler(layouts.Main()))

	fmt.Println("Server is running on port 3000")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}
