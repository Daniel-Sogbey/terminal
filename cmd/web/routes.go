package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(
		cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		},
	))

	mux.Use(middleware.Logger)
	mux.Use(NoSurf)

	mux.Get("/", Home)
	mux.Post("/initiate-payment", InitiatePayment)
	mux.Get("/paystack/callback", VerifyPayment)
	mux.Get("/success", Success)
	mux.Get("/failure", Failure)

	return mux
}
