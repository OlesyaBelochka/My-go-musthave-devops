package reporters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func sendUpdateRequestJson(fullPuth string, client http.Client, userData variables.Metrics) {

	strJSON, err := json.Marshal(userData)

	fmt.Println(string(strJSON))

	if err != nil {
		fmt.Errorf("marsalling failed: %v", err)
		return
	}

	resp, err := http.Post(fullPuth, "application/json", bytes.NewBuffer(strJSON))

	if err != nil {
		fmt.Errorf("post request failed: %v", err)
		return
	} else {
		fmt.Println(resp)
		defer resp.Body.Close()
	}
	return
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
	_, err := client.Do(req)
	if err != nil {
		log.Print("Sending failed", err)
		os.Exit(1)
	}

}

func Report(URL string, client http.Client) {

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

		sendUpdateRequestJson(URL, client, str)
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

		sendUpdateRequestJson(URL, client, str)
		variables.MC["PollCount"] = 0 // обнуляем?
	}

}
