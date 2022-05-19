package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"github.com/go-chi/chi/v5"
)

func sendStatus(w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status) // 404

	if status != http.StatusOK {
		http.Error(w, strconv.Itoa(status), status)
	}

}

func sendResponceJSON(w http.ResponseWriter, status int, needCompression bool, e string) {
	resp := variables.ServResponses{}

	if status == http.StatusOK {

		resp = variables.ServResponses{
			Result: "Data update succesfully",
			Error:  "",
		}
	} else {
		resp = variables.ServResponses{
			Result: "Unsuccesfully",
			Error:  e,
		}
	}

	strJSON, err := json.Marshal(resp)

	fmt.Println("#(sendResponceJSON) " + string(strJSON))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("(sendResponceJSON) Marshal error: " + err.Error()))
		fmt.Println("# (sendResponceJSON) Marshal error: " + err.Error())
		return
	}

	if needCompression {
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (sendResponceJSON) Compress error : " + err.Error()))
			fmt.Println("# (sendResponceJSON) Compress error : " + err.Error())
			return
		}
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(strJSON); err != nil {
		variables.PrinterErr(err, "HandleUpdateMetricsJSON"+"- Send error")
		return
	}

	fmt.Println("#(sendResponceJSON) responce succesfully")

}

func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {

	if variables.ShowLog {
		log.Print("HandleGetAllMetrics")
	}
	html := ""
	for s, c := range variables.MG {
		html += fmt.Sprintf("%s : %.3f\n", s, c)
	}

	for s, c := range variables.MC {

		html += fmt.Sprintf("%s : %d\n", s, c)

	}

	w.Header().Set("Content-Type", "text/html")

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		fmt.Println("(HandleGetAllMetrics)  сжимает файл чтобы отправить ответ")

		data, err := compression.Compress([]byte(html))
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		if _, err := w.Write(data); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := w.Write([]byte(html)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

}

func getMetric(mType, mName string, format bool) (string, int, error) {

	var answer string
	var st int
	var err error

	if mName == "" || mType == "" {
		st = http.StatusBadRequest
		return answer, st, err
	}

	switch strings.ToLower(mType) {
	case "gauge":

		if value, inMap := variables.MG[mName]; inMap {

			if format {
				//fmt.Println("getMetric with format")
				answer = fmt.Sprintf("%0.3f", value)
			} else {
				answer = strconv.FormatFloat(float64(value), 'f', -1, 64)

			}

			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}

	case "counter":
		if value, inMap := variables.MC[mName]; inMap {
			if format {
				//fmt.Println("getMetric with format")
				answer = fmt.Sprintf("%d", value)
			} else {
				//fmt.Println("getMetric except format")
				answer = strconv.FormatInt(int64(value), 10)
				//answer = fmt.Sprintf("%d", value)
			}
			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}
	default:

		st = http.StatusBadRequest
	}

	return answer, st, err
}

func HandleGetMetric(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("star HandleGetMetric old")

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	//fmt.Println("type metric: ", mType, " name metric: ", mName)

	val, code, err := getMetric(mType, mName, true)

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	fmt.Println("подобрали по типу: ", mType, " и  имени : ", mName, " значение метрики ", val, " в HandleGetMetric old")
	_, err = w.Write([]byte(val))

	if err != nil {
		fmt.Println(err)

	}

}

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
		fmt.Println(val)
		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		UpdateGaugeMetric(mName, variables.Gauge(val))

		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(mVal)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		UpdateCountMetric(mName, variables.Counter(val))

		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}

func readBodyJSONRequest(w http.ResponseWriter, r *http.Request, resp *variables.Metrics, needCompression *bool) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		*needCompression = true
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	if *needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#HandleGetMetricJSON mistake decompression: ")

	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("can't unmarshal: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// HandleGetMetricJSON возвращает метрику в виде JSON
func HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {

	var (
		resp            variables.Metrics
		needCompression bool
	)

	readBodyJSONRequest(w, r, &resp, &needCompression)

	mType := resp.MType
	mName := resp.ID

	val, code, err := getMetric(mType, mName, false)

	if err != nil {
		fmt.Println("#mistake getMetric: ", code, err)
		http.Error(w, err.Error(), code)
		return
	}
	switch mType {

	case "gauge":
		valFl, err := strconv.ParseFloat(val, 64)

		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Value = &valFl

	case "counter":
		valInt, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Delta = &valInt

	}

	strJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		fmt.Println("(HandleGetMetricJSON)  сжимает файл чтобы отправить ответ")
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (sendResponceJSON) Compress error : " + err.Error()))
			fmt.Println("# (sendResponceJSON) Compress error : " + err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(strJSON)

	if err != nil {
		fmt.Println(err)
	}
}

func HandleUpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {

	var (
		metrics         variables.Metrics
		needCompression bool
	)

	readBodyJSONRequest(w, r, &metrics, &needCompression)

	mType := metrics.MType
	mName := metrics.ID
	err := ""
	st := http.StatusNotImplemented

	if mName == "" || (mType != "gauge" && mType != "counter") {
		mType = ""
		err = "can't find gauge or counter or empty id"

	}

	switch strings.ToLower(mType) {

	case "gauge":
		val := *metrics.Value

		UpdateGaugeMetric(mName, variables.Gauge(val))
		st = http.StatusOK

	case "counter":
		val := *metrics.Delta

		UpdateCountMetric(mName, variables.Counter(val))
		st = http.StatusOK
	}

	sendResponceJSON(w, st, needCompression, err)

}

func UpdateGaugeMetric(name string, val variables.Gauge) {

	variables.MG[name] = val

}

func UpdateCountMetric(name string, val variables.Counter) {

	variables.MC[name] += val
}
