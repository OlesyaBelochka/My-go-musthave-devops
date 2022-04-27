package main

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
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

func sendRequest(fullPuth string, client http.Client) {

	data := url.Values{}
	req, _ := http.NewRequest(http.MethodPost, fullPuth, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(req)

	if err != nil {
		//os.Exit(1)
		panic(err)
	}

}

func getRequest(endpoint string, client http.Client) {

	var fullPuth string

	for k, v := range variables.MG {

		fullPuth = endpoint + strings.ToLower(strings.Replace(reflect.TypeOf(variables.MG[k]).String(), "variables.", "", -1)) + "/" + k + "/" + strconv.FormatFloat(float64(v), 'f', -1, 64)
		//log.Println(fullPuth)
		sendRequest(fullPuth, client)
	}

	for k, v := range variables.MC {
		fullPuth = endpoint + strings.ToLower(strings.Replace(reflect.TypeOf(variables.MC[k]).String(), "variables.", "", -1)) + "/" + k + "/" + strconv.FormatInt(int64(v), 10)
		//log.Println(fullPuth)
		sendRequest(fullPuth, client)

		variables.MC["PollCount"] = 0 // обнуляем????
	}

}
func main() {

	//fmt.Println("Начало...")
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
			//	fmt.Println("обновляем")
			runtime.ReadMemStats(st)
			updater.UpdateMetrics(st)

			timer10 = time.NewTimer(reportInterval * time.Second)
			//fmt.Println("отправляем")
			getRequest(endpoint, client)

		case <-timer.C:
			//fmt.Println("обновляем")
			runtime.ReadMemStats(st)

			updater.UpdateMetrics(st)

		case <-osSigChan:
			os.Exit(1)
			return

		}
	}

}
