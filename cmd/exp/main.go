package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/khongtrunght/lenslocked/controllers"
	"github.com/khongtrunght/lenslocked/templates/pages"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", controllers.StaticHandler(pages.Home()))
	router.Get("/about", controllers.StaticHandler(pages.About()))
	router.Get("/contact", controllers.StaticHandler(pages.Contact()))
	router.Get("/faq", controllers.StaticHandler(pages.Faq()))

	fmt.Println("Server is running on port 3000")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}
