package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

var mg = map[string]gauge{
	"Alloc":         0,
	"BuckHashSys":   0,
	"Frees":         0,
	"GCCPUFraction": 0,
	"GCSys":         0,
	"HeapAlloc":     0,
	"HeapIdle":      0,
	"HeapInuse":     0,
	"HeapObjects":   0,
	"HeapReleased":  0,
	"HeapSys":       0,
	"LastGC":        0,
	"MCacheInuse":   0,
	"MCacheSys":     0,
	"MSpanInuse":    0,
	"MSpanSys":      0,
	"Mallocs":       0,
	"NextGC":        0,
	"NumForcedGC":   0,
	"NumGC":         0,
	"RandomValue":   0,
}

var mc = map[string]counter{
	"PollCount": 0,
}

func sendStatusNotFound(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)

	http.Error(w, "Status not found", http.StatusNotFound)

}
func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(mg)
	json.NewEncoder(w).Encode(mc)
}

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {

	fmt.Println("HandleUpdateMetrics")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	var a = strings.Split(r.URL.String(), "/")
	var answer string
	if len(a) == 4 && (a[2] == "gauge" || a[2] == "counter") {

		if a[2] == "gauge" {
			if value, inMap := mg[a[3]]; inMap {
				answer = strconv.FormatFloat(float64(value), 'f', 10, 64)

			} else {

				sendStatusNotFound(w)
				return
			}

		} else {
			if value, inMap := mc[a[3]]; inMap {
				answer = strconv.FormatInt(int64(value), 10)
			} else {

				sendStatusNotFound(w)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(answer))

	} else {

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Wrong URL", http.StatusBadRequest)

	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {

	var a = strings.Split(r.URL.String(), "/")

	if len(a) == 5 && (a[2] == "gauge" || a[2] == "counter") {
		if a[2] == "gauge" {
			val2, _ := strconv.ParseFloat(a[4], 64)
			mg[a[3]] = gauge(val2)
		} else {
			val, _ := strconv.Atoi(a[4])
			mc["PollCount"] += counter(val)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		//	w.Write([]byte("метки обновились."))
	} else {

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Wrong URL", http.StatusNotFound)
	}

}
