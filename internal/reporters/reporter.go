package reporters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func sendUpdateRequestJSON(fullPuth string, client http.Client, userData variables.Metrics) {

	strJSON, err := json.Marshal(userData)
	//variables.FShowLog(string(strJSON))
	variables.PrinterErr(err, "")

	strJSON, err = compression.Compress(strJSON)

	// fmt.Println(strJSON)
	if err != nil {
		// вызываем останову агента и разбираемся. потом удали

		variables.PrinterErr(err, "# mitake during compression:")
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", fullPuth, bytes.NewBuffer(strJSON))
	variables.PrinterErr(err, "# mistake NewRequest : ")

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Content-Encoding", "gzip")

	resp, err := client.Do(req)
	variables.PrinterErr(err, "# sending mistake :")
	fmt.Println("клиент отправил: ", req)
	//resp, err := http.Post(fullPuth, "application/json", bytes.NewBuffer(strJSON))

	if resp != nil {
		err = resp.Body.Close()
		variables.PrinterErr(err, "")
	}

	//if resp.StatusCode != 200 {
	//	_, err := io.ReadAll(resp.Body)
	//	if err != nil {
	//		fmt.Errorf(err.Error())
	//		return
	//	}
	//	defer resp.Body.Close()
	//	//
	//}

}

func sendUpdateRequest(fullPuth string, client http.Client) {

	data := url.Values{}

	req, _ := http.NewRequest(http.MethodPost, fullPuth, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "text/plain")
	resp, err := client.Do(req)

	variables.PrinterErr(err, "")

	defer resp.Body.Close()

}

func gatgerData(client http.Client, URL string) {

	for k, v := range variables.MG {

		vFl := float64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &vFl,
		}
		//sendRequest(fmt.Sprintf("%sgauge/%s/%f", URL, k, v), client)
		if variables.ShowLog {
			log.Printf("отправляем метрику,  тип: %s , имя: %s, значение: %f", "gauge  в процедуре sendUpdateRequestJson", k, vFl)
		}

		sendUpdateRequestJSON(URL, client, str)
	}

	for k, v := range variables.MC {
		vInt := int64(v)
		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &vInt,
		}

		//sendRequest(fmt.Sprintf("%scounter/%s/%d", URL, k, v), client)

		if variables.ShowLog {
			log.Printf("отправляем метрику,  тип: %s , имя: %s, значение: %v", "counter", k, vInt)
		}

		sendUpdateRequestJSON(URL, client, str)
		variables.MC["PollCount"] = 0 // обнуляем?
	}
}

func Report(ctx context.Context, URL string) {

	client := http.Client{}

	for {
		timerReport := time.NewTimer(time.Duration(variables.ConfA.ReportInterval) * time.Second)

		select {
		case <-timerReport.C:
			// variables.FShowLog("sending...")
			gatgerData(client, URL)
		case <-ctx.Done():
			variables.FShowLog("ctx.Done(): Report")
			return

		}

	}
}
