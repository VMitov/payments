package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type key int

const (
	keyPayments key = iota
)

type api struct {
	db *sqlx.DB
}

func newAPI(dbconn string) (*api, error) {
	db, err := sqlx.Connect("postgres", dbconn)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to DB failed")
	}
	return &api{db: db}, nil
}

func newRouter(api *api) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/payments", func(r chi.Router) {
		r.Get("/", api.listPayments)
		r.Post("/", api.createPayment)

		r.Route("/{paymentID}", func(r chi.Router) {
			r.Get("/", api.getPayment)
		})

	})

	return r
}
