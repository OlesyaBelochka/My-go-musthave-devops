package updater

import (
	"math/rand"
	"runtime"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

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

	//variables.MG["MSpanInuse"] = variables.Gauge(st.MSpanInuse)
	//variables.MG["MSpanSys"] = variables.Gauge(st.MSpanSys)
	//variables.MG["Mallocs"] = variables.Gauge(st.Mallocs)
	//variables.MG["NextGC"] = variables.Gauge(st.NextGC)
	//variables.MG["NumForcedGC"] = variables.Gauge(st.NumForcedGC)
	//variables.MG["NumGC"] = variables.Gauge(st.NumGC)
	//variables.MG["OtherSys"] = variables.Gauge(st.OtherSys)
	//variables.MG["PauseTotalNs"] = variables.Gauge(st.PauseTotalNs)
	//variables.MG["StackInuse"] = variables.Gauge(st.StackInuse)
	//variables.MG["StackSys"] = variables.Gauge(st.StackSys)
	//variables.MG["Sys"] = variables.Gauge(st.Sys)
	//variables.MG["TotalAlloc"] = variables.Gauge(st.TotalAlloc)
	//variables.MG["RandomValue"] = variables.Gauge(rand.Int())

	//variables.MG["Alloc"] = variables.Gauge(st.Alloc)
	//variables.MG["BuckHashSys"] = variables.Gauge(st.BuckHashSys)
	//variables.MG["Frees"] = variables.Gauge(st.Frees)
	//variables.MG["GCCPUFraction"] = variables.Gauge(st.GCCPUFraction)
	//variables.MG["GCSys"] = variables.Gauge(st.GCSys)
	//variables.MG["HeapAlloc"] = variables.Gauge(st.HeapAlloc)
	//variables.MG["HeapIdle"] = variables.Gauge(st.HeapIdle)
	//variables.MG["HeapInuse"] = variables.Gauge(st.HeapInuse)
	//variables.MG["HeapObjects"] = variables.Gauge(st.HeapObjects)
	//variables.MG["HeapReleased"] = variables.Gauge(st.HeapReleased)
	//variables.MG["HeapSys"] = variables.Gauge(st.HeapSys)
	//variables.MG["LastGC"] = variables.Gauge(st.LastGC)
	//variables.MG["Lookups"] = variables.Gauge(st.Lookups)
	//variables.MG["MCacheInuse"] = variables.Gauge(st.MCacheInuse)
	//variables.MG["MCacheSys"] = variables.Gauge(st.MCacheSys)

	v := variables.MC["PollCount"] + 1
	UpdateCountMetric("PollCount", v)

}

func UpdateGaugeMetric(name string, val variables.Gauge) {

	variables.MG[name] = val

}

func UpdateCountMetric(name string, val variables.Counter) {

	variables.MC[name] = val

}
