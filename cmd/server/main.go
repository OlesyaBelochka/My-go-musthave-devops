package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	mux := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", HandleGetAllMetrics)
	mux.Get("/value/{mType}/{mName}", HandleGetMetric)
	mux.Get("/value/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {

		sendStatus(w, 700)

	})

	//mux.Get("/{anystring}", func(w http.ResponseWriter, r *http.Request) {
	//	sendStatus(w, 600)
	//	fmt.Println("anystring")
	//})

	mux.Post("/update/{mType}/{mName}/{mValue}", HandleUpdateMetrics)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))

}
