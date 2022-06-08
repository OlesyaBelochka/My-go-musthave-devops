package handlers

import (
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	val, code, err := getMetric(mType, mName, true)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	_, err = w.Write([]byte(val))

	if err != nil {
		variables.PrinterErr(err, "(HandleGetMetric) Ошибка отправки : ")
	}

}
