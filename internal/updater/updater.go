package updater

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

func Pall() {

	if variables.ShowLog {
		fmt.Println("#update..")
	}

	runtime.ReadMemStats(variables.MemSt)
	UpdateAllMetrics(variables.MemSt)

}

func UpdateAllMetrics(st *runtime.MemStats) {

	UpdateGaugeMetric("Alloc", variables.Gauge(st.Alloc))
	UpdateGaugeMetric("BuckHashSys", variables.Gauge(st.BuckHashSys))
	UpdateGaugeMetric("Frees", variables.Gauge(st.Frees))
	UpdateGaugeMetric("GCCPUFraction", variables.Gauge(st.GCCPUFraction))
	UpdateGaugeMetric("GCSys", variables.Gauge(st.GCSys))
	UpdateGaugeMetric("HeapAlloc", variables.Gauge(st.HeapAlloc))
	UpdateGaugeMetric("HeapIdle", variables.Gauge(st.HeapIdle))
	UpdateGaugeMetric("HeapObjects", variables.Gauge(st.HeapObjects))
	UpdateGaugeMetric("HeapReleased", variables.Gauge(st.HeapReleased))
	UpdateGaugeMetric("HeapSys", variables.Gauge(st.HeapSys))
	UpdateGaugeMetric("LastGC", variables.Gauge(st.LastGC))
	UpdateGaugeMetric("Lookups", variables.Gauge(st.Lookups))
	UpdateGaugeMetric("MCacheInuse", variables.Gauge(st.MCacheInuse))
	UpdateGaugeMetric("MCacheSys", variables.Gauge(st.MCacheSys))
	UpdateGaugeMetric("MSpanInuse", variables.Gauge(st.MSpanInuse))
	UpdateGaugeMetric("MSpanSys", variables.Gauge(st.MSpanSys))
	UpdateGaugeMetric("Mallocs", variables.Gauge(st.Mallocs))
	UpdateGaugeMetric("NextGC", variables.Gauge(st.NextGC))
	UpdateGaugeMetric("NumForcedGC", variables.Gauge(st.NumForcedGC))
	UpdateGaugeMetric("NumGC", variables.Gauge(st.NumGC))
	UpdateGaugeMetric("OtherSys", variables.Gauge(st.OtherSys))
	UpdateGaugeMetric("PauseTotalNs", variables.Gauge(st.PauseTotalNs))
	UpdateGaugeMetric("StackInuse", variables.Gauge(st.StackInuse))
	UpdateGaugeMetric("StackSys", variables.Gauge(st.StackSys))
	UpdateGaugeMetric("Sys", variables.Gauge(st.Sys))
	UpdateGaugeMetric("TotalAlloc", variables.Gauge(st.TotalAlloc))
	UpdateGaugeMetric("RandomValue", variables.Gauge(rand.Int()))

	v := variables.MC["PollCount"] + 1
	UpdateCountMetric("PollCount", v)

}

func UpdateGaugeMetric(name string, val variables.Gauge) {

	if variables.ShowLog {
		log.Printf("обновляем метку %v  в значение %v", name, val)
	}

	variables.MG[name] = val

}

func UpdateCountMetric(name string, val variables.Counter) {

	if variables.ShowLog {
		log.Printf("обновляем метку %v  в значение %v", name, val)
	}

	variables.MC[name] = val

}
