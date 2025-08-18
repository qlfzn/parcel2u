package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qlfzn/parcel2u/internal/auth"
)

// application struct
type Application struct {
	Addr string
}

func (a *Application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
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
