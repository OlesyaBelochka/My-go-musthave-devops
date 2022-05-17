package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/updater"
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
	w.Header().Set("Accept-Encoding", "gzip")

	_, err = w.Write(strJSON)

	if err != nil {
		variables.PrinterErr(err, "HandleUpdateMetricsJSON"+"- Send error")
		return
	}

	fmt.Println("#(sendResponceJSON) responce succesfully")

}

func HandleGetAllMetrics(w http.ResponseWriter, r *http.Request) {

	if variables.ShowLog {
		log.Print("HandleGetAllMetrics")
	}

	w.Header().Set("Content-Type", "text/plain")

	html := ""
	for s, c := range variables.MG {

		html += fmt.Sprintf("%s : %.3f\n", s, c)

	}
	for s, c := range variables.MC {

		html += fmt.Sprintf("%s : %d\n", s, c)
	}
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(html))

	if err != nil {
		fmt.Println(err)
	}
}

func getMetric(mType, mName string, format bool) (string, int, error) {

	var answer string
	var st int
	var err error

	//if variables.ShowLog {
	//	fmt.Println(mName)
	//	fmt.Println(mType)
	//}

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

	//if variables.ShowFullLog {
	//fmt.Println("вот такой ответ дала процедура getMetric", answer, st, err)
	//}
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

func HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {

	fmt.Println("star HandleGetMetric Json")

	var (
		//metrics         variables.Metrics
		needCompression bool
	)

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			needCompression = true

		}
	}

	if needCompression {

		fmt.Print("(HandleGetMetricJSON) из агента пришли данные о том, что нужна компрессия, ", r.Header.Get("Accept-Encoding"), r.Header.Get("Content-Encoding"))

	} else {

		fmt.Print("(HandleGetMetricJSON) из агента пришли данные о том, что НЕ  нужна комрессия , ")
	}

	var resp variables.Metrics

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	if needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#HandleGetMetricJSON mistake decompression: ")
	}

	//fmt.Println("GetMetricJSON Handler: " + string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(w, "can't unmarshal: ", err.Error())
		return
	}

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

		if variables.ShowLog {
			fmt.Println("#подобрали по типу: ", mType, " и  имени : ", mName, " значение метрики ", valFl, " в HandleGetMetric Json")
		}

	case "counter":
		valInt, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			fmt.Println(err)
			return
		}
		resp.Delta = &valInt

		if variables.ShowLog {
			fmt.Println("#подобрали по типу: ", mType, " и  имени : ", mName, " значение метрики ", valInt, " в HandleGetMetric Json")
		}

	}

	strJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	//f

	if needCompression {

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

	fmt.Println("#GetMetricJSON Handler: "+string(body), " answer ", string(strJSON))

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept-Encoding", "gzip")

	_, err = w.Write(strJSON)

	if err != nil {
		fmt.Println(err)
	}

}

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("HandleUpdateMetrics old")
	//var a = strings.Split(r.URL.String(), "/")

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mVal := chi.URLParam(r, "mValue")

	if variables.ShowLog {
		fmt.Println("получен URL ", r.URL)
		fmt.Printf("распознали: mType=%s, mName=%s, mVal = %v", mType, mName, mVal)
	}

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

		updater.UpdateGaugeMetric(mName, variables.Gauge(val))

		sendStatus(w, http.StatusOK)

	case "counter":

		val, err := strconv.Atoi(mVal)

		if err != nil {
			sendStatus(w, http.StatusBadRequest) // 400
			return
		}

		updater.UpdateCountMetric(mName, variables.Counter(val))

		sendStatus(w, http.StatusOK)

	default:
		sendStatus(w, http.StatusNotImplemented) // 501
	}

}

func HandleUpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {

	//var a = strings.Split(r.URL.String(), "/")
	var (
		metrics         variables.Metrics
		needCompression bool
	)

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			needCompression = true

		}
	}

	if needCompression {

		fmt.Print("(HandleUpdateMetricsJSON) из агента пришли данные о том, что нужна компрессия, ", r.Header.Get("Accept-Encoding"), r.Header.Get("Content-Encoding"))

	} else {

		fmt.Print("(HandleUpdateMetricsJSON) из агента пришли данные о том, что НЕ  нужна комрессия , ")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#UpdateMetricsJSON mistake decompression: ")
	}

	fmt.Println("#UpdateMetricsJSON Handler: " + string(body))
	err = json.Unmarshal(body, &metrics)

	if err != nil {
		fmt.Println(w, "can't unmarshal: ", err.Error())
	}

	mType := metrics.MType
	mName := metrics.ID

	if mName == "" || (mType != "gauge" && mType != "counter") {

		sendResponceJSON(w, http.StatusNotImplemented, false, "can't find gauge or counter or empty id") // 501
		return
	}

	switch strings.ToLower(mType) {

	case "gauge":
		val := *metrics.Value

		updater.UpdateGaugeMetric(mName, variables.Gauge(val))
		sendResponceJSON(w, http.StatusOK, needCompression, "")

	case "counter":

		val := *metrics.Delta

		updater.UpdateCountMetric(mName, variables.Counter(val))
		sendResponceJSON(w, http.StatusOK, needCompression, "")

	default:
		sendResponceJSON(w, http.StatusNotImplemented, needCompression, "default switch") // 501
	}

}
