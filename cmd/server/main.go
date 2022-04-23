package main

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/handlers"
	"log"
	"net/http"
)

func main() {

	//http.HandleFunc("/update/", handlers.HandleMetrics)
	//log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))

	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/update/", handlers.HandleMetrics)
	//mux.HandleFunc("/update/counter/", handlers.HandleGaugeC)
	//
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))

}
