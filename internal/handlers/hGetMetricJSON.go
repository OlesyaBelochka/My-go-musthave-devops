package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/prhash"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"strconv"
	"strings"
)

// HandleGetMetricJSON возвращает метрику в виде JSON
func HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {

	var (
		resp            variables.Metrics
		needCompression bool
	)

	st, err := readBodyJSONRequest(w, r, &resp, &needCompression, false)

	if err != nil {
		// это значит что мы по какой-то причине не смогли выполнить процедуру выше и отправляем статус
		sendStatus(w, st) // 400
		return
	}

	if internal.ConfS.Key != "" {
		variables.FShowLog(fmt.Sprintf("при отправке с сервера вычислять хеш = %v", internal.ConfS.Key))
	}

	mType := resp.MType
	mName := resp.ID

	val, code, err := getMetric(mType, mName, false)

	if err != nil {
		variables.PrinterErr(err, "(HandleGetMetricJSON) ошибка выполнения функции getMetric: ")
		http.Error(w, err.Error(), code)
		return
	}

	switch mType {

	case "gauge":
		valFl, err := strconv.ParseFloat(val, 64)

		if err != nil {
			variables.PrinterErr(err, "(HandleGetMetricJSON) не можем сделать ParseFloat для gauge")
			return
		}
		resp.Value = &valFl
		if internal.ConfS.Key != "" {
			resp.Hash = prhash.Hash(fmt.Sprintf("%s:gauge:%f", resp.ID, valFl), internal.ConfS.Key)
		}

	case "counter":
		valInt, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			variables.PrinterErr(err, "(HandleGetMetricJSON) не можем сделать ParseInt для counter")
			return
		}

		resp.Delta = &valInt
		if internal.ConfS.Key != "" {
			resp.Hash = prhash.Hash(fmt.Sprintf("%s:counter:%d", resp.ID, valInt), internal.ConfS.Key)
		}
	}

	strJSON, err := json.Marshal(resp)
	if err != nil {
		variables.PrinterErr(err, "(HandleGetMetricJSON) не можем сделать json.Marshal для структуры")
		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		variables.FShowLog("(HandleGetMetricJSON) Увидели заголовок Accept-Encoding,  сжимает файл чтобы отправить ответ ")
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (HandleGetMetricJSON) Compress error : " + err.Error()))
			variables.PrinterErr(err, "# (HandleGetMetricJSON) Compress error : ")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if internal.ConfS.Key != "" {
		variables.FShowLog("отправили структуру на сервер с хэшами так как сервер имеет ключ" + string(strJSON))
	}

	_, err = w.Write(strJSON)

	if err != nil {
		variables.PrinterErr(err, "#(HandleGetMetricJSON) ошибка отправки ответного файла от сервера: ")
	}
}
