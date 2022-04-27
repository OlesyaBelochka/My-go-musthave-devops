package main

import (
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

func sendRequest(fullPuth string, client http.Client) {

	if variables.ShowLog {
		fmt.Println(fullPuth)
	}

	data := url.Values{}
	req, _ := http.NewRequest(http.MethodPost, fullPuth, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(req)

	if err != nil {
		os.Exit(1)
		log.Printf("Sending failed", err)
	}

}

func getRequest(URL string, client http.Client) {

	for k, v := range variables.MG {

		sendRequest(fmt.Sprintf("%sgauge/%s/%f", URL, k, v), client)
	}

	for k, v := range variables.MC {

		sendRequest(fmt.Sprintf("%scounter/%s/%d", URL, k, v), client)

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
