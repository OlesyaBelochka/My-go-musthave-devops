package handlers

import (
	"context"
	"errors"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleUpdatesSliceMetricsJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	var (
		slMetrics       []variables.Metrics
		needCompression bool

		mmNameG []string
		mmNameC []string
		mmValG  [][]byte
		mmValC  [][]byte
	)

	st, err := readBodySliceJSONRequest(w, r, &slMetrics, &needCompression)

	if err == nil {

		for _, metrics := range slMetrics {

			mType := metrics.MType
			mName := metrics.ID

			if mName == "" || (mType != "gauge" && mType != "counter") {
				mType = ""
				err = errors.New("can't find gauge or counter or empty id")
				st = http.StatusNotImplemented
			}

			switch strings.ToLower(mType) {

			case "gauge":

				val := *metrics.Value

				err := verificationHash(metrics, config.VarConfServer.Key)

				if err != nil {
					variables.PrinterErr(err, "")
					return
				}
				mmNameG = append(mmNameG, mName)
				mmValG = append(mmValG, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

			case "counter":

				val := *metrics.Delta

				err := verificationHash(metrics, config.VarConfServer.Key)
				if err != nil {
					variables.PrinterErr(err, "")
					return
				}
				mmNameC = append(mmNameC, mName)
				mmValC = append(mmValC, []byte(strconv.FormatInt(int64(val), 10)))

			}
		}

		if len(mmNameG) > 0 {
			storage.MGServer.SetSlice(ctx, mmNameG, mmValG)
			st = http.StatusOK
		}

		if len(mmNameC) > 0 {
			storage.MCServer.SetSlice(ctx, mmNameC, mmValC)
			st = http.StatusOK
		}

		sendResponceJSON(w, st, needCompression, err)

	}
}
