package inmemory

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"math/rand"
	"runtime"
	"strconv"
)

type (
	MCounter map[string]variables.Counter
	MGauge   map[string]variables.Gauge
)

type MemoryStorage struct {
	MC MCounter
	MG MGauge
}

func New() *MemoryStorage {
	return &MemoryStorage{
		MC: MCounter{},
		MG: MGauge{},
	}

}

func (M MGauge) Set(name string, val []byte) {
	byteToFloat, _ := strconv.ParseFloat(string(val), 64)
	M[name] = variables.Gauge(byteToFloat)
	//fmt.Printf("Set Gauge %s, in val = %f \n", name, byteToFloat)
}

func (M MCounter) Set(name string, val []byte) {
	byteToInt, _ := strconv.ParseInt(string(val), 10, 64)
	M[name] += variables.Counter(byteToInt)
	//fmt.Printf("Set Counter %s, in val = %d \n", name, M[name])
}

func (M MGauge) Get(s string, val []byte) {
	//TODO implement me
	//	panic("implement me")
}
func (M MCounter) Get(s string, val []byte) {
	//TODO implement me
	//panic("implement me")
}

func (M MGauge) Pall(st *runtime.MemStats) {
	runtime.ReadMemStats(variables.MemSt)

	M.Set("Alloc", []byte(strconv.FormatFloat(float64(st.Alloc), 'f', -1, 64)))
	M.Set("BuckHashSys", []byte(strconv.FormatFloat(float64(st.BuckHashSys), 'f', -1, 64)))
	M.Set("Frees", []byte(strconv.FormatFloat(float64(st.Frees), 'f', -1, 64)))
	M.Set("GCCPUFraction", []byte(strconv.FormatFloat(float64(st.GCCPUFraction), 'f', -1, 64)))
	M.Set("GCSys", []byte(strconv.FormatFloat(float64(st.GCSys), 'f', -1, 64)))
	M.Set("HeapAlloc", []byte(strconv.FormatFloat(float64(st.HeapAlloc), 'f', -1, 64)))
	M.Set("HeapIdle", []byte(strconv.FormatFloat(float64(st.HeapIdle), 'f', -1, 64)))
	M.Set("HeapInuse", []byte(strconv.FormatFloat(float64(st.HeapInuse), 'f', -1, 64)))
	M.Set("HeapObjects", []byte(strconv.FormatFloat(float64(st.HeapObjects), 'f', -1, 64)))
	M.Set("HeapReleased", []byte(strconv.FormatFloat(float64(st.HeapReleased), 'f', -1, 64)))
	M.Set("HeapSys", []byte(strconv.FormatFloat(float64(st.HeapSys), 'f', -1, 64)))
	M.Set("LastGC", []byte(strconv.FormatFloat(float64(st.LastGC), 'f', -1, 64)))
	M.Set("Lookups", []byte(strconv.FormatFloat(float64(st.Lookups), 'f', -1, 64)))
	M.Set("MCacheInuse", []byte(strconv.FormatFloat(float64(st.MCacheInuse), 'f', -1, 64)))
	M.Set("MCacheSys", []byte(strconv.FormatFloat(float64(st.MCacheSys), 'f', -1, 64)))
	M.Set("MSpanInuse", []byte(strconv.FormatFloat(float64(st.MSpanInuse), 'f', -1, 64)))
	M.Set("MSpanSys", []byte(strconv.FormatFloat(float64(st.MSpanSys), 'f', -1, 64)))
	M.Set("Mallocs", []byte(strconv.FormatFloat(float64(st.Mallocs), 'f', -1, 64)))
	M.Set("NextGC", []byte(strconv.FormatFloat(float64(st.NextGC), 'f', -1, 64)))
	M.Set("NumForcedGC", []byte(strconv.FormatFloat(float64(st.NumForcedGC), 'f', -1, 64)))
	M.Set("NumGC", []byte(strconv.FormatFloat(float64(st.NumGC), 'f', -1, 64)))
	M.Set("OtherSys", []byte(strconv.FormatFloat(float64(st.OtherSys), 'f', -1, 64)))
	M.Set("PauseTotalNs", []byte(strconv.FormatFloat(float64(st.PauseTotalNs), 'f', -1, 64)))
	M.Set("StackInuse", []byte(strconv.FormatFloat(float64(st.StackInuse), 'f', -1, 64)))
	M.Set("StackSys", []byte(strconv.FormatFloat(float64(st.StackSys), 'f', -1, 64)))
	M.Set("Sys", []byte(strconv.FormatFloat(float64(st.Sys), 'f', -1, 64)))
	M.Set("TotalAlloc", []byte(strconv.FormatFloat(float64(st.TotalAlloc), 'f', -1, 64)))
	M.Set("RandomValue", []byte(strconv.FormatFloat(float64(rand.Int()), 'f', -1, 64)))

}

func (M MCounter) Pall(st *runtime.MemStats) {

	M.Set("PollCount", []byte(strconv.FormatInt(int64(1), 10)))
}
