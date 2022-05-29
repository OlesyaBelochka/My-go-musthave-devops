package reporters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/prhash"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"time"
)

type ReporterInterface interface {
	Report(key string)
}

type gaugeM struct {
	name string
	val  float64
}

type counterM struct {
	name string
	val  int64
}

type GaugeReporter struct {
	M []gaugeM
}

type CounterReporter struct {
	M []counterM
}

func NewGaugeReporter() *GaugeReporter {

	g := make([]gaugeM, 0)
	for k, v := range storage.MGAgent.M {

		gM := gaugeM{name: k, val: float64(v)}
		g = append(g, gM)
	}
	return &GaugeReporter{M: g}
}

func NewCounterReporter() *CounterReporter {

	g := make([]counterM, 0)

	for k, v := range storage.MCAgent.M {
		gC := counterM{name: k, val: int64(v)}
		g = append(g, gC)

	}

	return &CounterReporter{M: g}
}

func (g GaugeReporter) Report(key string) {
	var str []variables.Metrics

	for _, v := range g.M {
		v := v
		str1 := variables.Metrics{
			ID:    v.name,
			MType: "gauge",
			Value: &v.val,
			Hash:  prhash.Hash(fmt.Sprintf("%s:gauge:%f", v.name, v.val), key),
		}

		str = append(str, str1)

	}

	SendButchJSON(str)
}

func (g CounterReporter) Report(key string) {
	var str []variables.Metrics

	for _, v := range g.M {
		v := v
		str1 := variables.Metrics{
			ID:    v.name,
			MType: "counter",
			Delta: &v.val,
			Hash:  prhash.Hash(fmt.Sprintf("%s:counter:%d", v.name, v.val), key),
		}
		str = append(str, str1)
		//SendJSON(str1)
	}
	SendButchJSON(str)
}

func SendButchJSON(userData []variables.Metrics) {

	if len(userData) == 0 {
		fmt.Println("Агент получил пустую структуру. На сервер ее не обрабатываем, и не отправляем")
		return
	}
	strJSON, err := json.MarshalIndent(userData, "", "  ")
	fmt.Println("такую структуру отправляем на сервер: ", string(strJSON))

	variables.PrinterErr(err, "")
	strJSON, err = compression.Compress(strJSON)
	variables.PrinterErr(err, "# mitake during compression:")

	req, err := http.NewRequest("POST", "http://"+internal.ConfA.Address+"/updates/", bytes.NewBuffer(strJSON))
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

func ReportAgentNew(ctx context.Context, key string) {

	for {
		timerReport := time.NewTimer(internal.ConfA.ReportInterval)
		select {
		case <-timerReport.C:
			gR := NewGaugeReporter()
			cR := NewCounterReporter()
			fmt.Println(gR)
			fmt.Println(cR)

			var reporters []ReporterInterface
			rep := append(reporters, gR, cR)

			for _, reporterInterface := range rep {
				reporterInterface.Report(key)
			}

			variables.FShowLog("#reporting..")
			//		reportMetrics()
		case <-ctx.Done():
			variables.FShowLog("ctx.Done(): Report")
			return

		}

	}
}
