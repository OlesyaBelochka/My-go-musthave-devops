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

var AgentMetrics = inmemory.New()
var ServerMetrics = inmemory.New()

func PallMetrics(ch Storage) {

	runtime.ReadMemStats(variables.MemSt)
	ch.Pall(variables.MemSt)

}
