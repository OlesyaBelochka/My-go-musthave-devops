package storage

import (
	"context"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/inmemory"
)

type Storage interface {
	Set(s string, val []byte)
	Get(s string) ([]byte, bool)
	SetSlice(ctx context.Context, s []string, val [][]byte)
}

var MGAgent = inmemory.NewGaugeMS()
var MCAgent = inmemory.NewCounterMS()

var MGServer Storage
var MCServer Storage
