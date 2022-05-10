package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/joho/godotenv"
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

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func sendUpdateRequestJson(fullPuth string, client http.Client, userData *variables.Metrics) {

	strJSON, err := json.MarshalIndent(*userData, "", "	")

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
		if variables.ShowLog {
			log.Printf("отправляем метрику,  тип: %s , имя: %s, значение: %f", "gauge  в процедуре sendUpdateRequestJson", k, v_fl)
		}
		sendUpdateRequestJson(URL, client, &str)
	}

	for k, v := range variables.MC {
		v_int := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &v_int,
		}

		//sendRequest(fmt.Sprintf("%scounter/%s/%d", URL, k, v), client)

		if variables.ShowLog {
			log.Printf("отправляем метрику,  тип: %s , имя: %s, значение: %v", "counter", k, v_int)
		}

		sendUpdateRequestJson(URL, client, &str)
		variables.MC["PollCount"] = 0 // обнуляем?
	}

}
func main() {
	log.Println("Клиент запустился, обновляет и отправляет")
	conf := config.New()

	if variables.ShowLog {
		fmt.Printf("Address %v, ReportInterval = %v, PollInterval =  %v", conf.Address, conf.ReportInterval, conf.PollInterval)
	}

	st := new(runtime.MemStats)

	endpoint := "/update/"

	client := http.Client{}

	timer10 := time.NewTimer(time.Duration(conf.ReportInterval) * time.Second)

	for {
		osSigChan := make(chan os.Signal)
		signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		timer := time.NewTimer(time.Duration(conf.PollInterval) * time.Second)

		select {
		case <-timer10.C:

			if variables.ShowLog {
				fmt.Println("#update..")
			}

			runtime.ReadMemStats(st)
			updater.UpdateAllMetrics(st)

			timer10 = time.NewTimer(time.Duration(conf.ReportInterval) * time.Second)

			if variables.ShowLog {
				fmt.Println("#send to ", conf.Address+endpoint)
			}

			getRequest("http://"+conf.Address+endpoint, client)

		case <-timer.C:

			if variables.ShowLog {
				fmt.Println("#update..")
			}

			runtime.ReadMemStats(st)

			updater.UpdateAllMetrics(st)

		case <-osSigChan:
			os.Exit(1)
			return
		}
	}

}
