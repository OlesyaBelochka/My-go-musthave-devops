package main

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if variables.ShowLog {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Get("/", handlers.HandleGetAllMetrics)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{mType}/{mName}", handlers.HandleGetMetric)
	})

	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdateMetrics)

	http.ListenAndServe(variables.IPServer, r)

}
