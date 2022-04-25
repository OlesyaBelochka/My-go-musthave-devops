package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	//mux := chi.NewMux()
	//
	//// зададим встроенные middleware, чтобы улучшить стабильность приложения
	//mux.Use(middleware.RequestID)
	//mux.Use(middleware.RealIP)
	//mux.Use(middleware.Logger)
	//mux.Use(middleware.Recoverer)
	//
	//mux.Get("/", HandleGetAllMetrics)
	//mux.Get("/value/{mType}/{mName}", HandleGetMetric)
	//mux.Get("/value/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {
	//
	//	sendStatus(w, 700)
	//
	//})
	//
	////mux.Get("/{anystring}", func(w http.ResponseWriter, r *http.Request) {
	////	sendStatus(w, 600)
	////	fmt.Println("anystring")
	////})
	//
	//mux.Post("/update/{mType}/{mName}/{mValue}", HandleUpdateMetrics)
	////localhost:8080
	//log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
	////log.Fatal(http.ListenAndServe(":8080", mux))

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
		// GET /value/Gauge
		r.Get("/{mType}", HandleGetMetric)
		r.Get("/{mType}/", HandleGetMetric)
		// GET /value/Gauge/GCCPUFraction
		r.Get("/{mType}/{mName}", HandleGetMetric)
		r.Get("/{mType}/{mName}/", HandleGetMetric)

		//
		// GETGET /value/Gauge/GCCPUFraction/1
		//r.Get("/{mType}/{mName}/{mValue}/", HandleGetMetric)
	})

	r.Post("/update/{mType}/{mName}/{mValue}", HandleUpdateMetrics)

	http.ListenAndServe("127.0.0.1:8080", r)

	//http.ListenAndServe(":8080", r)
	//address := "127.0.0.1:8080"
	//
	//log.Fatal(http.ListenAndServe(address, mux))

}
