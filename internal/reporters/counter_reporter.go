package reporters

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/prhash"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

type counterM struct {
	name string
	val  int64
}

type CounterReporter struct {
	M []counterM
}

func NewCounterReporter() *CounterReporter {

	g := make([]counterM, 0)

	for k, v := range storage.MCAgent.M {
		gC := counterM{name: k, val: int64(v)}
		g = append(g, gC)

	}

	return &CounterReporter{M: g}
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
