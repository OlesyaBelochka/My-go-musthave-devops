package reporters

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/prhash"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

type gaugeM struct {
	name string
	val  float64
}

type GaugeReporter struct {
	M []gaugeM
}

func NewGaugeReporter() *GaugeReporter {

	g := make([]gaugeM, 0)
	for k, v := range storage.MGAgent.M {

		gM := gaugeM{name: k, val: float64(v)}
		g = append(g, gM)
	}
	return &GaugeReporter{M: g}
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
