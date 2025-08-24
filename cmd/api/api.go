package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/qlfzn/parcel2u/internal/auth"
)

// application struct
type Application struct {
	Addr        string
	AuthHandler *auth.Handler
}

func (a *Application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, 
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/users", a.AuthHandler.RegisterUser)
		r.Post("/login", a.AuthHandler.LoginUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("you are authorised!"))
		})
	})

	return r
}

func (a *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    a.Addr,
		Handler: mux,
	}

	log.Printf("server has started at: http://localhost%s\n", a.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
