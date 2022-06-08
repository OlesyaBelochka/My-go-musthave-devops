package handlers

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
)

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mVal := chi.URLParam(r, "mValue")

	if mName == "" || mVal == "" || (mType != "gauge" && mType != "counter") {
		sendStatus(w, http.StatusNotImplemented) // 501
		return
	}

	switch strings.ToLower(mType) {

	case "gauge":
		val, err := strconv.ParseFloat(mVal, 64)
		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		storage.MGServer.Set(mName, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(mVal)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		storage.MCServer.Set(mName, []byte(strconv.FormatInt(int64(val), 10)))

		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}
