package main

import (
	"fmt"
	"net/http"

	"companyintersala/internal/api"
	dbpkg "companyintersala/internal/api/db"
	authMiddleware "companyintersala/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbpkg.InitDB("./users.db")
	defer dbpkg.DB.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/signup", api.SignupHandler)
	r.Post("/login", api.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware)
		r.Get("/profile", api.ProfileHandler)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.AdminOnly)
			r.Get("/users", api.UsersHandler)
		})
	})

	fmt.Println("server starting on :8080")
	http.ListenAndServe(":8080", r)

}
