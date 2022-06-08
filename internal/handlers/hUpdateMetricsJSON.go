package handlers

import (
	"errors"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"strconv"
	"strings"
)

func HandleUpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {

	var (
		metrics         variables.Metrics
		needCompression bool
	)

	st, err := readBodyJSONRequest(w, r, &metrics, &needCompression, true)

	if err == nil {

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

			storage.MGServer.Set(mName, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

			st = http.StatusOK

		case "counter":
			val := *metrics.Delta

			storage.MCServer.Set(mName, []byte(strconv.FormatInt(int64(val), 10)))

			st = http.StatusOK
		}
	}

	sendResponceJSON(w, st, needCompression, err)

}
