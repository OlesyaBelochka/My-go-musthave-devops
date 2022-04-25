package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", HandleGetAllMetrics)
	r.Route("/value", func(r chi.Router) {
		// GET /value
		r.Get("/", HandleGetMetric)

		r.Get("/{mType}", HandleGetMetric)

	})

	r.Post("/update/{mType}/{mName}/{mValue}", HandleUpdateMetrics)

	http.ListenAndServe("127.0.0.1:8080", r)

}
