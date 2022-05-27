package inmemory

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"strconv"
)

type (
	MCounter map[string]variables.Counter
	MGauge   map[string]variables.Gauge
)

type CounterMemoryStorage struct {
	M MCounter
}

type GaugeMemoryStorage struct {
	M MGauge
}

func NewCounterMS() *CounterMemoryStorage {
	return &CounterMemoryStorage{
		M: MCounter{},
	}
}

func NewGaugeMS() *GaugeMemoryStorage {
	return &GaugeMemoryStorage{
		M: MGauge{},
	}
}

func (M GaugeMemoryStorage) Set(name string, val []byte) {
	byteToFloat, _ := strconv.ParseFloat(string(val), 64)
	M.M[name] = variables.Gauge(byteToFloat)
	fmt.Printf("Set Gauge %s, in val = %f \n", name, byteToFloat)
}

func (M CounterMemoryStorage) Set(name string, val []byte) {
	byteToInt, _ := strconv.ParseInt(string(val), 10, 64)
	M.M[name] += variables.Counter(byteToInt)
	fmt.Printf("Set Counter %s, in val = %d \n", name, M.M[name])
}

func (M GaugeMemoryStorage) Get(s string) ([]byte, bool) {

	if value, inMap := M.M[s]; inMap {
		return []byte(strconv.FormatFloat(float64(value), 'f', -1, 64)), true
	}
	return []byte(""), false // пустой список байт

}
func (M CounterMemoryStorage) Get(s string) ([]byte, bool) {

	if value, inMap := M.M[s]; inMap {
		return []byte(strconv.FormatInt(int64(value), 10)), true
	}
	return []byte(""), false // пустой список байт
}
