package poller

import (
	"context"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"time"
)

func PallStart(ctx context.Context) {
	for {
		timer := time.NewTimer(internal.ConfA.PollInterval)
		select {
		case <-timer.C:
			variables.FShowLog("#update..")
			storage.PallMetrics(storage.AgentMetrics.MG)
			storage.PallMetrics(storage.AgentMetrics.MC)

		case <-ctx.Done():
			variables.FShowLog("ctx.Done()")
			return
		}
	}

}
