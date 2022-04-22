package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type gauge float64
type counter int64

var metGauge gauge
var metCounter counter

func HandleGaugeC(w http.ResponseWriter, r *http.Request) {
	var val int64
	fmt.Println("counte")

	u, _ := url.Parse(r.URL.Redacted())
	println("Path: ", u.Path)
	val, err := strconv.ParseInt(path.Base(u.Path), 10, 64)

	if err != nil {
		http.Error(w, "userId is empty", http.StatusBadRequest)
	}

	metCounter = metCounter + counter(val)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
func HandleGaugeM(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Gauge")
	fmt.Println(r.URL.Redacted())
	fmt.Println(r.URL.Query())
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {

	http.HandleFunc("/update/", HandleMetrics)
	http.ListenAndServe("127.0.0.1:8080", nil)

	//mux := http.NewServeMux()
	//
	//mux.Handle("/update/", http.NotFoundHandler())
	//mux.HandleFunc("/update/gauge/", HandleGaugeM)
	//mux.HandleFunc("/update/counter/", HandleGaugeC)
	//
	//log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))

	/*
		mux := http.NewServeMux()

		// Create sample handler to returns 404
		mux.Handle("/update", http.NotFoundHandler())

		// Create sample handler that returns 200
		mux.Handle("/update/counter/", newPeopleHandler())

		log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
	*/
}
