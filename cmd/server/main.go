package main

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	mux := chi.NewRouter()
	fmt.Println(variables.MC["Alloc"])
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", handlers.HandleGetAllMetrics)
	mux.Get("/update/{mType}/{mName}", handlers.HandleGetMetric)
	mux.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdateMetrics)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}
