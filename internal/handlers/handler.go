package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

var mm = map[string]gauge{
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

//func HandleGaugeC(w http.ResponseWriter, r *http.Request) {
//	var val int64
//
//	u, _ := url.Parse(r.URL.Redacted())
//	println("Path: ", u.Path)
//	val, err := strconv.ParseInt(path.Base(u.Path), 10, 64)
//
//	if err != nil {
//		http.Error(w, "Wrong URL", http.StatusBadRequest)
//	}
//
//	w.Header().Set("Content-Type", "text/plain")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(`{"status":"ok"}`))
//}

//func HandleGaugeM(w http.ResponseWriter, r *http.Request) {
//
//	fmt.Println("Gauge")
//	fmt.Println(r.URL.Redacted())
//	fmt.Println(r.URL.Query())
//
//	var a = strings.Split(r.URL.String(), "/")
//	fmt.Println(len(a))
//
//	for i, s := range a {
//		fmt.Println(i, "=", s)
//	}
//
//	if len(a) == 5 && a[2] == "gauge" {
//		val2, _ := strconv.ParseFloat(a[4], 64)
//		mm[a[3]] = gauge(val2)
//	} else {
//		val, _ := strconv.Atoi(a[4])
//		mc["PollCount"] += counter(val)
//	}
//
//	w.Header().Set("Content-Type", "text/plain")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(`{"status":"ok"}`))
//}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("Gauge")
	//fmt.Println(r.URL.Redacted())
	//fmt.Println(r.URL.Query())

	var a = strings.Split(r.URL.String(), "/")
	//fmt.Println(len(a))

	if len(a) == 5 && (a[2] == "gauge" || a[2] == "counter") {
		if a[2] == "gauge" {
			val2, _ := strconv.ParseFloat(a[4], 64)
			mm[a[3]] = gauge(val2)
		} else {
			val, _ := strconv.Atoi(a[4])
			mc["PollCount"] += counter(val)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		//w.Write()
	} else {

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Wrong URL", http.StatusBadRequest)

	}

}
