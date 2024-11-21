package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/khongtrunght/lenslocked/controllers"
	"github.com/khongtrunght/lenslocked/models"
	"github.com/khongtrunght/lenslocked/templates"
	"github.com/khongtrunght/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	cfg := models.DefaultPostgresConfig()

	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signin.gohtml"))
	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Post("/users", usersC.Create)
	r.Get("/users/me", usersC.CurrentUser)

	csrfKey := "csrf_token"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false)) // TODO: set to true in production

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", csrfMw(r))
}
