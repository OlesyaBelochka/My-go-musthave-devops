package inmemory

import (
	"context"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"strconv"
	"sync"
)

type (
	MCounter map[string]variables.Counter
	MGauge   map[string]variables.Gauge
)

type CounterMemoryStorage struct {
	M   MCounter
	Mtx sync.RWMutex
}

type GaugeMemoryStorage struct {
	M   MGauge
	Mtx sync.RWMutex
}

func NewCounterMS() *CounterMemoryStorage {
	return &CounterMemoryStorage{
		M:   MCounter{},
		Mtx: sync.RWMutex{},
	}
}

func NewGaugeMS() *GaugeMemoryStorage {
	return &GaugeMemoryStorage{
		M:   MGauge{},
		Mtx: sync.RWMutex{},
	}
}

func (M *GaugeMemoryStorage) Set(name string, val []byte) {

	byteToFloat, _ := strconv.ParseFloat(string(val), 64)
	//M.Mtx.Lock()
	M.M[name] = variables.Gauge(byteToFloat)
	//M.Mtx.Unlock()
	variables.FShowLog(fmt.Sprintf("(Set :GaugeMemoryStorage)  %s, in val = %f \n", name, byteToFloat))
}

func (M *CounterMemoryStorage) Set(name string, val []byte) {

	byteToInt, _ := strconv.ParseInt(string(val), 10, 64)
	//M.Mtx.Lock()
	M.M[name] += variables.Counter(byteToInt)
	//M.Mtx.Unlock()
	variables.FShowLog(fmt.Sprintf("(Set: CounterMemoryStorage) %s, in val = %d \n", name, M.M[name]))
}

func (M *GaugeMemoryStorage) SetSlice(ctx context.Context, name []string, val [][]byte) {
	for i := 0; i < len(name); i++ {
		M.Set(name[i], val[i])
	}
}

func (M *CounterMemoryStorage) SetSlice(ctx context.Context, name []string, val [][]byte) {
	for i := 0; i < len(name); i++ {
		M.Set(name[i], val[i])
	}
}

func (M *GaugeMemoryStorage) Get(s string) ([]byte, bool) {
	//M.Mtx.RLock()
	//defer M.Mtx.RUnlock()

	if value, inMap := M.M[s]; inMap {
		return []byte(strconv.FormatFloat(float64(value), 'f', -1, 64)), true
	}
	return []byte(""), false // пустой список байт

}
func (M *CounterMemoryStorage) Get(s string) ([]byte, bool) {
	//M.Mtx.RLock()
	//defer M.Mtx.RUnlock()

	if value, inMap := M.M[s]; inMap {
		return []byte(strconv.FormatInt(int64(value), 10)), true
	}
	return []byte(""), false // пустой список байт
}
