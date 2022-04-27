package updater

import (
	"math/rand"
	"runtime"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

func UpdateMetrics(st *runtime.MemStats) {

	variables.MG["Alloc"] = variables.Gauge(st.Alloc)
	variables.MG["BuckHashSys"] = variables.Gauge(st.BuckHashSys)
	variables.MG["Frees"] = variables.Gauge(st.Frees)
	variables.MG["GCCPUFraction"] = variables.Gauge(st.GCCPUFraction)
	variables.MG["GCSys"] = variables.Gauge(st.GCSys)
	variables.MG["HeapAlloc"] = variables.Gauge(st.HeapAlloc)
	variables.MG["HeapIdle"] = variables.Gauge(st.HeapIdle)
	variables.MG["HeapInuse"] = variables.Gauge(st.HeapInuse)
	variables.MG["HeapObjects"] = variables.Gauge(st.HeapObjects)
	variables.MG["HeapReleased"] = variables.Gauge(st.HeapReleased)
	variables.MG["HeapSys"] = variables.Gauge(st.HeapSys)
	variables.MG["LastGC"] = variables.Gauge(st.LastGC)
	variables.MG["Lookups"] = variables.Gauge(st.Lookups)
	variables.MG["MCacheInuse"] = variables.Gauge(st.MCacheInuse)
	variables.MG["MCacheSys"] = variables.Gauge(st.MCacheSys)
	variables.MG["MSpanInuse"] = variables.Gauge(st.MSpanInuse)
	variables.MG["MSpanSys"] = variables.Gauge(st.MSpanSys)
	variables.MG["Mallocs"] = variables.Gauge(st.Mallocs)
	variables.MG["NextGC"] = variables.Gauge(st.NextGC)
	variables.MG["NumForcedGC"] = variables.Gauge(st.NumForcedGC)
	variables.MG["NumGC"] = variables.Gauge(st.NumGC)
	variables.MG["OtherSys"] = variables.Gauge(st.OtherSys)
	variables.MG["PauseTotalNs"] = variables.Gauge(st.PauseTotalNs)
	variables.MG["StackInuse"] = variables.Gauge(st.StackInuse)
	variables.MG["StackSys"] = variables.Gauge(st.StackSys)
	variables.MG["Sys"] = variables.Gauge(st.Sys)
	variables.MG["TotalAlloc"] = variables.Gauge(st.TotalAlloc)
	variables.MG["RandomValue"] = variables.Gauge(rand.Int())

	variables.MC["PollCount"]++
}
