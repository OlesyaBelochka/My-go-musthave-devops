package inmemory

import (
	"context"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"strconv"
	"sync"
)

type (
	MGauge map[string]variables.Gauge
)

type GaugeMemoryStorage struct {
	M   MGauge
	Mtx sync.RWMutex
}

func NewGaugeMS() *GaugeMemoryStorage {
	return &GaugeMemoryStorage{
		M:   MGauge{},
		Mtx: sync.RWMutex{},
	}
}

func (M *GaugeMemoryStorage) Set(name string, val []byte) {

	byteToFloat, _ := strconv.ParseFloat(string(val), 64)
	M.Mtx.Lock()
	M.M[name] = variables.Gauge(byteToFloat)
	M.Mtx.Unlock()
	variables.FShowLog(fmt.Sprintf("(Set :GaugeMemoryStorage)  %s, in val = %f \n", name, byteToFloat))
}

func (M *GaugeMemoryStorage) SetSlice(ctx context.Context, name []string, val [][]byte) {
	for i := 0; i < len(name); i++ {
		M.Set(name[i], val[i])
	}
}

func (M *GaugeMemoryStorage) Get(s string) ([]byte, bool) {

	if value, inMap := M.M[s]; inMap {
		return []byte(strconv.FormatFloat(float64(value), 'f', -1, 64)), true
	}
	return []byte(""), false // пустой список байт

}
