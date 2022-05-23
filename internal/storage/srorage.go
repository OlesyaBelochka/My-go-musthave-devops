package storage

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/inmemory"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"runtime"
)

type Storage interface {
	Set(s string, val []byte)
	Get(s string, val []byte)
	Pall(st *runtime.MemStats)
}

func PallMetrics(ch Storage) {
	runtime.ReadMemStats(variables.MemSt)
	ch.Pall(variables.MemSt)
}

var MGAgent = inmemory.NewGaugeMS()
var MCAgent = inmemory.NewCounterMS()

var MGServer = inmemory.NewGaugeMS()
var MCServer = inmemory.NewCounterMS()
