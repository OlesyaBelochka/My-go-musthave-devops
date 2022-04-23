package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"strconv"
	"strings"
)

func sendStatusNotFound(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)

	http.Error(w, "Status not found", http.StatusNotFound)

}
func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleUpdateMetrics")

	json.NewEncoder(w).Encode(variables.MG)
	json.NewEncoder(w).Encode(variables.MC)
}

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {

	fmt.Println("HandleGetMetric")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	var a = strings.Split(r.URL.String(), "/")
	var answer string
	if len(a) == 4 && (strings.ToLower(a[2]) == "gauge" || strings.ToLower(a[2]) == "counter") {

		if strings.ToLower(a[2]) == "gauge" {
			if value, inMap := variables.MG[a[3]]; inMap {
				answer = strconv.FormatFloat(float64(value), 'f', 10, 64)

			} else {

				sendStatusNotFound(w)
				return
			}

		} else {
			if value, inMap := variables.MC[a[3]]; inMap {
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
	fmt.Println("HandleUpdateMetrics")
	var a = strings.Split(r.URL.String(), "/")

	if len(a) == 5 && (strings.ToLower(a[2]) == "gauge" || strings.ToLower(a[2]) == "counter") {
		if strings.ToLower(a[2]) == "gauge" {
			val2, _ := strconv.ParseFloat(a[4], 64)
			variables.MG[a[3]] = variables.Gauge(val2)
		} else {
			val, _ := strconv.Atoi(a[4])
			variables.MC["PollCount"] += variables.Counter(val)
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
