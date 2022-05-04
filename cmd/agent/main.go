package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

const pollInterval = 2
const reportInterval = 10

func sendUpdateRequestJson(fullPuth string, client http.Client, userData variables.Metrics) {

	strJSON, err := json.MarshalIndent(userData, "", "	")

	if err != nil {
		fmt.Errorf("marsalling failed: %v", err)
	}

	_, err = http.Post(fullPuth, "application/json", bytes.NewBuffer(strJSON))

	if err != nil {
		log.Print("Sending failed", err)
		os.Exit(1)
	}

}

func sendUpdateRequest(fullPuth string, client http.Client) {

	data := url.Values{}

	req, _ := http.NewRequest(http.MethodPost, fullPuth, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(req)
	if err != nil {
		log.Print("Sending failed", err)
		os.Exit(1)
	}

}

func getRequest(URL string, client http.Client) {

	for k, v := range variables.MG {

		v_fl := float64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &v_fl,
		}
		//sendRequest(fmt.Sprintf("%sgauge/%s/%f", URL, k, v), client)
		sendUpdateRequestJson(URL, client, str)
	}

	for k, v := range variables.MC {
		v_int := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Delta: &v_int,
		}

		//sendRequest(fmt.Sprintf("%scounter/%s/%d", URL, k, v), client)
		sendUpdateRequestJson(URL, client, str)
		variables.MC["PollCount"] = 0 // обнуляем?
	}

}
func main() {
	st := new(runtime.MemStats)

	endpoint := "/update/"

	client := http.Client{}

	timer10 := time.NewTimer(reportInterval * time.Second)

	for {
		osSigChan := make(chan os.Signal)
		signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		timer := time.NewTimer(pollInterval * time.Second)

		select {
		case <-timer10.C:

			if variables.ShowLog {
				fmt.Println("#update..")
			}

			runtime.ReadMemStats(st)
			updater.UpdateMetrics(st)

			timer10 = time.NewTimer(reportInterval * time.Second)

			if variables.ShowLog {
				fmt.Println("#send..")
			}

			getRequest("http://"+variables.IPServer+endpoint, client)

		case <-timer.C:

			if variables.ShowLog {
				fmt.Println("#update..")
			}

			runtime.ReadMemStats(st)

			updater.UpdateMetrics(st)

		case <-osSigChan:
			os.Exit(1)
			return
		}
	}

}
