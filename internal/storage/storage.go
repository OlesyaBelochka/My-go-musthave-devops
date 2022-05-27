package storage

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/inmemory"
)

type Storage interface {
	Set(s string, val []byte)
	Get(s string) ([]byte, bool)
	//Pall(st *runtime.MemStats)
}

var MGAgent = inmemory.NewGaugeMS()
var MCAgent = inmemory.NewCounterMS()

var MGServer Storage
var MCServer Storage
