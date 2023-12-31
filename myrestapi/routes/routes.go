package routes

import (
	"net/http"

	"github.com/Daniel-Sogbey/myrestapi/handlers"
	"github.com/Daniel-Sogbey/myrestapi/helpers"
	"github.com/Daniel-Sogbey/myrestapi/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userDBStore := models.NewUserDBStore(models.DB)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		payload := &models.DataResponse{
			Status:  "success",
			Message: "Server live",
			Data:    "Hello, World",
		}

		helpers.WriteJSON(w, payload, http.StatusOK)
	})

	mux.Post("/signup", handlers.Signup)
	// mux.Get("/login", handlers.Login)

	// mux.With(middleware.Auth).Patch("/forgot-password", handlers.ResetPassword)

	return mux
}

// 2 timothy 3:14-17
