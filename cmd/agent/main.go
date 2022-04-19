package main

import (
	"fmt"
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

type gauge float64
type counter int64

type MyMetrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge

	RandomValue gauge
	PollCount   counter
}

const pollInterval = 2
const reportInterval = 10

func (mm *MyMetrics) upDateMetrics(st *runtime.MemStats) {

	mm.Alloc = gauge(st.Alloc)
	mm.BuckHashSys = gauge(st.BuckHashSys)
	mm.Frees = gauge(st.Frees)
	mm.GCCPUFraction = gauge(st.GCCPUFraction)
	mm.GCSys = gauge(st.GCSys)
	mm.HeapAlloc = gauge(st.HeapAlloc)
	mm.HeapIdle = gauge(st.HeapIdle)
	mm.HeapInuse = gauge(st.HeapInuse)
	mm.HeapObjects = gauge(st.HeapObjects)
	mm.HeapReleased = gauge(st.HeapReleased)
	mm.HeapSys = gauge(st.HeapSys)
	mm.LastGC = gauge(st.LastGC)
	mm.Lookups = gauge(st.Lookups)
	mm.MCacheInuse = gauge(st.MCacheInuse)
	mm.MCacheSys = gauge(st.MCacheSys)
	mm.MSpanInuse = gauge(st.MSpanInuse)
	mm.MSpanSys = gauge(st.MSpanSys)
	mm.Mallocs = gauge(st.Mallocs)
	mm.NextGC = gauge(st.NextGC)
	mm.NumForcedGC = gauge(st.NumForcedGC)
	mm.NumGC = gauge(st.NumGC)
	mm.OtherSys = gauge(st.OtherSys)
	mm.PauseTotalNs = gauge(st.PauseTotalNs)
	mm.StackInuse = gauge(st.StackInuse)
	mm.StackSys = gauge(st.StackSys)
	mm.Sys = gauge(st.Sys)
	mm.TotalAlloc = gauge(st.TotalAlloc)
	mm.RandomValue = gauge(rand.Int())
	mm.PollCount++

}

func sendRequest(fullPuth string, client http.Client) {
	data := url.Values{}
	req, _ := http.NewRequest(http.MethodPost, fullPuth, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(req)

	if err != nil {
		panic(err)
	}

}

func getRequest(endpoint string, client http.Client, mm MyMetrics) {
	var fullPuth string

	//mm.Alloc = gauge(st.Alloc)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.Alloc).String(), "main.", "", -1) + "/" + "Alloc" + "/" + strconv.FormatFloat(float64(mm.Alloc), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.BuckHashSys = gauge(st.BuckHashSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.BuckHashSys).String(), "main.", "", -1) + "/" + "BuckHashSys" + "/" + strconv.FormatFloat(float64(mm.BuckHashSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.Frees = gauge(st.Frees)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.Frees).String(), "main.", "", -1) + "/" + "Frees" + "/" + strconv.FormatFloat(float64(mm.Frees), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.GCCPUFraction = gauge(st.GCCPUFraction)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.GCCPUFraction).String(), "main.", "", -1) + "/" + "GCCPUFraction" + "/" + strconv.FormatFloat(float64(mm.GCCPUFraction), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//m.GCSys = gauge(st.GCSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.GCSys).String(), "main.", "", -1) + "/" + "GCSys" + "/" + strconv.FormatFloat(float64(mm.GCSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//	mm.HeapAlloc = gauge(st.HeapAlloc)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.HeapAlloc).String(), "main.", "", -1) + "/" + "HeapAlloc" + "/" + strconv.FormatFloat(float64(mm.HeapAlloc), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.HeapIdle = gauge(st.HeapIdle)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.HeapIdle).String(), "main.", "", -1) + "/" + "HeapIdle" + "/" + strconv.FormatFloat(float64(mm.HeapIdle), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.HeapInuse = gauge(st.HeapInuse)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.HeapInuse).String(), "main.", "", -1) + "/" + "HeapInuse" + "/" + strconv.FormatFloat(float64(mm.HeapInuse), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.HeapObjects = gauge(st.HeapObjects)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.HeapObjects).String(), "main.", "", -1) + "/" + "HeapObjects" + "/" + strconv.FormatFloat(float64(mm.HeapObjects), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.HeapReleased = gauge(st.HeapReleased)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.HeapReleased).String(), "main.", "", -1) + "/" + "HeapReleased" + "/" + strconv.FormatFloat(float64(mm.HeapReleased), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.HeapSys = gauge(st.HeapSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.HeapSys).String(), "main.", "", -1) + "/" + "HeapSys" + "/" + strconv.FormatFloat(float64(mm.HeapSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.LastGC = gauge(st.LastGC)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.LastGC).String(), "main.", "", -1) + "/" + "LastGC" + "/" + strconv.FormatFloat(float64(mm.LastGC), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.Lookups = gauge(st.Lookups)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.Lookups).String(), "main.", "", -1) + "/" + "Lookups" + "/" + strconv.FormatFloat(float64(mm.Lookups), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.MCacheInuse = gauge(st.MCacheInuse)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.MCacheInuse).String(), "main.", "", -1) + "/" + "MCacheInuse" + "/" + strconv.FormatFloat(float64(mm.MCacheInuse), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.MCacheSys = gauge(st.MCacheSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.MCacheSys).String(), "main.", "", -1) + "/" + "MCacheSys" + "/" + strconv.FormatFloat(float64(mm.MCacheSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.MSpanInuse = gauge(st.MSpanInuse)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.MSpanInuse).String(), "main.", "", -1) + "/" + "MSpanInuse" + "/" + strconv.FormatFloat(float64(mm.MSpanInuse), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.MSpanSys = gauge(st.MSpanSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.MSpanSys).String(), "main.", "", -1) + "/" + "MSpanSys" + "/" + strconv.FormatFloat(float64(mm.MSpanSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.Mallocs = gauge(st.Mallocs)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.Mallocs).String(), "main.", "", -1) + "/" + "Mallocs" + "/" + strconv.FormatFloat(float64(mm.Mallocs), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.NextGC = gauge(st.NextGC)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.NextGC).String(), "main.", "", -1) + "/" + "NextGC" + "/" + strconv.FormatFloat(float64(mm.NextGC), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.NumForcedGC = gauge(st.NumForcedGC)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.NumForcedGC).String(), "main.", "", -1) + "/" + "NumForcedGC" + "/" + strconv.FormatFloat(float64(mm.NumForcedGC), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.NumGC = gauge(st.NumGC)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.NumGC).String(), "main.", "", -1) + "/" + "NumGC" + "/" + strconv.FormatFloat(float64(mm.NumGC), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.OtherSys = gauge(st.OtherSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.OtherSys).String(), "main.", "", -1) + "/" + "OtherSys" + "/" + strconv.FormatFloat(float64(mm.OtherSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.PauseTotalNs = gauge(st.PauseTotalNs)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.PauseTotalNs).String(), "main.", "", -1) + "/" + "PauseTotalNs" + "/" + strconv.FormatFloat(float64(mm.PauseTotalNs), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.StackInuse = gauge(st.StackInuse)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.StackInuse).String(), "main.", "", -1) + "/" + "StackInuse" + "/" + strconv.FormatFloat(float64(mm.StackInuse), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.StackSys = gauge(st.StackSys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.StackSys).String(), "main.", "", -1) + "/" + "StackSys" + "/" + strconv.FormatFloat(float64(mm.StackSys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.Sys = gauge(st.Sys)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.Sys).String(), "main.", "", -1) + "/" + "Sys" + "/" + strconv.FormatFloat(float64(mm.Sys), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.TotalAlloc = gauge(st.TotalAlloc)
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.TotalAlloc).String(), "main.", "", -1) + "/" + "TotalAlloc" + "/" + strconv.FormatFloat(float64(mm.TotalAlloc), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.RandomValue = gauge(rand.Int())
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.RandomValue).String(), "main.", "", -1) + "/" + "RandomValue" + "/" + strconv.FormatFloat(float64(mm.RandomValue), 'f', -1, 64)
	sendRequest(fullPuth, client)

	//mm.PollCount++
	fullPuth = endpoint + strings.Replace(reflect.TypeOf(mm.PollCount).String(), "main.", "", -1) + "/" + "PollCount" + "/" + strconv.FormatInt(int64(mm.PollCount), 10)
	fmt.Println(fullPuth)
	sendRequest(fullPuth, client)

}
func main() {
	fmt.Println("Начало...1")
	st := new(runtime.MemStats)

	var mm MyMetrics

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
			mm.upDateMetrics(st)

			timer10 = time.NewTimer(reportInterval * time.Second)
			fmt.Println("отправляем")
			fmt.Println(mm.Alloc)
			fmt.Println(mm.BuckHashSys)
			fmt.Println(mm.Frees)
			fmt.Println(mm.GCCPUFraction)
			fmt.Println(mm.PollCount)
			fmt.Println(mm.RandomValue)
			fmt.Println("--------")
			//u, _ := url.ParseRequestURI(apiUrl)

			getRequest(endpoint, client, mm)

		case <-timer.C:
			fmt.Println("обновляем")
			runtime.ReadMemStats(st)
			mm.upDateMetrics(st)

		case <-osSigChan:
			break
		}
	}

}
