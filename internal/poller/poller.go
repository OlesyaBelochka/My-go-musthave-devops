package poller

import (
	"context"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func PallStart(ctx context.Context) {
	for {
		timer := time.NewTimer(config.VarConfAgent.PollInterval)
		select {
		case <-timer.C:

			variables.FShowLog("#update..")

			PallMetrics()

		case <-ctx.Done():
			variables.FShowLog("ctx.Done()")
			return
		}
	}

}

func pallRuntimeMetrics() {
	st := new(runtime.MemStats)
	runtime.ReadMemStats(st)

	storage.MGAgent.Set("Alloc", []byte(strconv.FormatFloat(float64(st.Alloc), 'f', -1, 64)))
	storage.MGAgent.Set("BuckHashSys", []byte(strconv.FormatFloat(float64(st.BuckHashSys), 'f', -1, 64)))
	storage.MGAgent.Set("Frees", []byte(strconv.FormatFloat(float64(st.Frees), 'f', -1, 64)))
	storage.MGAgent.Set("GCCPUFraction", []byte(strconv.FormatFloat(float64(st.GCCPUFraction), 'f', -1, 64)))
	storage.MGAgent.Set("GCSys", []byte(strconv.FormatFloat(float64(st.GCSys), 'f', -1, 64)))
	storage.MGAgent.Set("HeapAlloc", []byte(strconv.FormatFloat(float64(st.HeapAlloc), 'f', -1, 64)))
	storage.MGAgent.Set("HeapIdle", []byte(strconv.FormatFloat(float64(st.HeapIdle), 'f', -1, 64)))
	storage.MGAgent.Set("HeapInuse", []byte(strconv.FormatFloat(float64(st.HeapInuse), 'f', -1, 64)))
	storage.MGAgent.Set("HeapObjects", []byte(strconv.FormatFloat(float64(st.HeapObjects), 'f', -1, 64)))
	storage.MGAgent.Set("HeapReleased", []byte(strconv.FormatFloat(float64(st.HeapReleased), 'f', -1, 64)))
	storage.MGAgent.Set("HeapSys", []byte(strconv.FormatFloat(float64(st.HeapSys), 'f', -1, 64)))
	storage.MGAgent.Set("LastGC", []byte(strconv.FormatFloat(float64(st.LastGC), 'f', -1, 64)))
	storage.MGAgent.Set("Lookups", []byte(strconv.FormatFloat(float64(st.Lookups), 'f', -1, 64)))
	storage.MGAgent.Set("MCacheInuse", []byte(strconv.FormatFloat(float64(st.MCacheInuse), 'f', -1, 64)))
	storage.MGAgent.Set("MCacheSys", []byte(strconv.FormatFloat(float64(st.MCacheSys), 'f', -1, 64)))
	storage.MGAgent.Set("MSpanInuse", []byte(strconv.FormatFloat(float64(st.MSpanInuse), 'f', -1, 64)))
	storage.MGAgent.Set("MSpanSys", []byte(strconv.FormatFloat(float64(st.MSpanSys), 'f', -1, 64)))
	storage.MGAgent.Set("Mallocs", []byte(strconv.FormatFloat(float64(st.Mallocs), 'f', -1, 64)))
	storage.MGAgent.Set("NextGC", []byte(strconv.FormatFloat(float64(st.NextGC), 'f', -1, 64)))
	storage.MGAgent.Set("NumForcedGC", []byte(strconv.FormatFloat(float64(st.NumForcedGC), 'f', -1, 64)))
	storage.MGAgent.Set("NumGC", []byte(strconv.FormatFloat(float64(st.NumGC), 'f', -1, 64)))
	storage.MGAgent.Set("OtherSys", []byte(strconv.FormatFloat(float64(st.OtherSys), 'f', -1, 64)))
	storage.MGAgent.Set("PauseTotalNs", []byte(strconv.FormatFloat(float64(st.PauseTotalNs), 'f', -1, 64)))
	storage.MGAgent.Set("StackInuse", []byte(strconv.FormatFloat(float64(st.StackInuse), 'f', -1, 64)))
	storage.MGAgent.Set("StackSys", []byte(strconv.FormatFloat(float64(st.StackSys), 'f', -1, 64)))
	storage.MGAgent.Set("Sys", []byte(strconv.FormatFloat(float64(st.Sys), 'f', -1, 64)))
	storage.MGAgent.Set("TotalAlloc", []byte(strconv.FormatFloat(float64(st.TotalAlloc), 'f', -1, 64)))

	storage.MGAgent.Set("RandomValue", []byte(strconv.FormatFloat(float64(rand.Int()), 'f', -1, 64)))
	storage.MCAgent.Set("PollCount", []byte(strconv.FormatInt(int64(1), 10)))
}

func pallGopsutilMetrics() {

	m, err := mem.VirtualMemory()
	if err != nil {
		variables.PrinterErr(err, "не удалось получить метки памяти")
		return
	}
	storage.MGAgent.Set("TotalMemory", []byte(strconv.FormatFloat(float64(m.Total), 'f', -1, 64)))
	storage.MGAgent.Set("FreeMemory", []byte(strconv.FormatFloat(float64(m.Free), 'f', -1, 64)))

	percentage, _ := cpu.Percent(0, true)

	for idx, cpupercent := range percentage {
		storage.MGAgent.Set(fmt.Sprintf("CPUutilization%d", idx+1), []byte(strconv.FormatFloat(float64(cpupercent), 'f', -1, 64)))
	}

}

func PallMetrics() {

	pallRuntimeMetrics()
	pallGopsutilMetrics()

}
