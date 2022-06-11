package reporters

import (
	"bytes"
	"encoding/json"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
)

type Reporter interface {
	Report(key string)
}

func SendButchJSON(userData []variables.Metrics) {

	if len(userData) == 0 {
		variables.FShowLog("Агент получил пустую структуру. На сервер ее не обрабатываем, и не отправляем")
		return
	}
	strJSON, err := json.MarshalIndent(userData, "", "  ")
	variables.FShowLog("такую структуру отправляем на сервер: " + string(strJSON))
	variables.PrinterErr(err, "(SendButchJSON) Не смогли сделать Маршал ошибка: ")

	strJSON, err = compression.Compress(strJSON)
	variables.PrinterErr(err, "# mitake during compression:")

	req, err := http.NewRequest("POST", "http://"+config.VarConfAgent.Address+"/updates/", bytes.NewBuffer(strJSON))
	variables.PrinterErr(err, "# mistake NewRequest : ")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := config.Client.Do(req)
	variables.PrinterErr(err, "# sending mistake :")
	if resp != nil {
		err = resp.Body.Close()
		variables.PrinterErr(err, "")
	}
}

func SendJSON(userData variables.Metrics) { //использовалась ранее когда отправляли ответ не списком
	strJSON, err := json.Marshal(userData)
	variables.PrinterErr(err, "")
	strJSON, err = compression.Compress(strJSON)
	variables.PrinterErr(err, "# mitake during compression:")

	req, err := http.NewRequest("POST", "http://"+config.VarConfAgent.Address+config.EndpointAgent, bytes.NewBuffer(strJSON))
	variables.PrinterErr(err, "# mistake NewRequest : ")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := config.Client.Do(req)
	variables.PrinterErr(err, "# sending mistake :")
	if resp != nil {
		err = resp.Body.Close()
		variables.PrinterErr(err, "")
	}
}
