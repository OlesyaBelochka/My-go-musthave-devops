package reporters

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"time"
)

func reportMetrics() {
	reportGauge()
	reportCounter()
}

func reportGauge() {
	for k, v := range storage.AgentMetrics.MG {
		vFl := float64(v)

		str := variables.Metrics{
			ID:    k,
			MType: "gauge",
			Value: &vFl,
		}
		SendJSON(str)
	}
}

func reportCounter() {
	for k, v := range storage.AgentMetrics.MC {
		vInt := int64(v)

		str := variables.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &vInt,
		}
		SendJSON(str)
	}
}

func SendJSON(userData variables.Metrics) {
	strJSON, err := json.Marshal(userData)
	variables.PrinterErr(err, "")
	strJSON, err = compression.Compress(strJSON)
	variables.PrinterErr(err, "# mitake during compression:")

	req, err := http.NewRequest("POST", "http://"+internal.ConfA.Address+internal.EndpointAgent, bytes.NewBuffer(strJSON))
	variables.PrinterErr(err, "# mistake NewRequest : ")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := internal.Client.Do(req)
	variables.PrinterErr(err, "# sending mistake :")
	if resp != nil {
		err = resp.Body.Close()
		variables.PrinterErr(err, "")
	}
}

func ReportAgent(ctx context.Context) {

	for {
		timerReport := time.NewTimer(internal.ConfA.ReportInterval)
		select {
		case <-timerReport.C:
			variables.FShowLog("#reporting..")
			reportMetrics()
		case <-ctx.Done():
			variables.FShowLog("ctx.Done(): Report")
			return

		}

	}
}
