package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const pollInterval = 2
const reportInterval = 10

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

func upDateMetrics(st *runtime.MemStats) {

	mm["Alloc"] = gauge(st.Alloc)
	mm["BuckHashSys"] = gauge(st.BuckHashSys)
	mm["Frees"] = gauge(st.Frees)
	mm["GCCPUFraction"] = gauge(st.GCCPUFraction)
	mm["GCSys"] = gauge(st.GCSys)
	mm["HeapAlloc"] = gauge(st.HeapAlloc)
	mm["HeapIdle"] = gauge(st.HeapIdle)
	mm["HeapInuse"] = gauge(st.HeapInuse)
	mm["HeapObjects"] = gauge(st.HeapObjects)
	mm["HeapReleased"] = gauge(st.HeapReleased)
	mm["HeapSys"] = gauge(st.HeapSys)
	mm["LastGC"] = gauge(st.LastGC)
	mm["Lookups"] = gauge(st.Lookups)
	mm["MCacheInuse"] = gauge(st.MCacheInuse)
	mm["MCacheSys"] = gauge(st.MCacheSys)
	mm["MSpanInuse"] = gauge(st.MSpanInuse)
	mm["MSpanSys"] = gauge(st.MSpanSys)
	mm["Mallocs"] = gauge(st.Mallocs)
	mm["NextGC"] = gauge(st.NextGC)
	mm["NumForcedGC"] = gauge(st.NumForcedGC)
	mm["NumGC"] = gauge(st.NumGC)
	mm["OtherSys"] = gauge(st.OtherSys)
	mm["PauseTotalNs"] = gauge(st.PauseTotalNs)
	mm["StackInuse"] = gauge(st.StackInuse)
	mm["StackSys"] = gauge(st.StackSys)
	mm["Sys"] = gauge(st.Sys)
	mm["TotalAlloc"] = gauge(st.TotalAlloc)
	mm["RandomValue"] = gauge(rand.Int())

	mc["PollCount"]++

}
func sendRequest(fullPuth string, client http.Client) {
	data := url.Values{}
	req, _ := http.NewRequest(http.MethodPost, fullPuth, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func getRequest(endpoint string, client http.Client) {

	var fullPuth string

	for k, v := range mm {
		fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm[k]).String(), "main.", "", -1) + "/" + k + "/" + strconv.FormatFloat(float64(v), 'f', -1, 64)
		log.Println(fullPuth)
		sendRequest(fullPuth, client)
	}

	for k, v := range mc {
		fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm[k]).String(), "main.", "", -1) + "/" + k + "/" + strconv.FormatFloat(float64(v), 'f', -1, 64)
		log.Println(fullPuth)
		sendRequest(fullPuth, client)
		mc["PollCount"] = 0 // обнуляем
	}

}
func main() {

	fmt.Println("Начало...")
	st := new(runtime.MemStats)

	endpoint := "http://127.0.0.1:8080/update/"

	client := http.Client{}

	timer10 := time.NewTimer(reportInterval * time.Second)

	for {
		osSigChan := make(chan os.Signal)
		signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		timer := time.NewTimer(pollInterval * time.Second)

		select {
		case <-timer10.C:
			fmt.Println("обновляем")
			runtime.ReadMemStats(st)
			upDateMetrics(st)

			timer10 = time.NewTimer(reportInterval * time.Second)
			fmt.Println("отправляем")
			getRequest(endpoint, client)

		case <-timer.C:
			fmt.Println("обновляем")
			runtime.ReadMemStats(st)
			upDateMetrics(st)

		case <-osSigChan:
			break
		}
	}

}
