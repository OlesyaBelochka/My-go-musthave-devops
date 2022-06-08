package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/prhash"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func verificationHash(resp variables.Metrics, partKey2 string) error {

	if resp.Hash != "" {

		getSHash := ""
		partKey1 := ""
		switch resp.MType {

		case "gauge":
			partKey1 = fmt.Sprintf("%s:%s:%f", resp.ID, resp.MType, *resp.Value)
		case "counter":
			partKey1 = fmt.Sprintf("%s:%s:%d", resp.ID, resp.MType, *resp.Delta)
		default:
			return errors.New("не смогли посчитать хэш на сервере так как получен неверны тип метрики")
		}

		getSHash = prhash.Hash(partKey1, partKey2)

		str := fmt.Sprintf("полученный хеш: %s, посчитанный хеш: %s,считали для метки: %s,использовали ключ: %s", resp.Hash, getSHash, partKey1, partKey2)

		if getSHash != resp.Hash {
			return errors.New("Хеши не равны: " + str)
		}
	}

	return nil

}

func sendStatus(w http.ResponseWriter, status int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status) // 404

	if status != http.StatusOK {
		http.Error(w, strconv.Itoa(status), status)
	}

}

func sendResponceJSON(w http.ResponseWriter, status int, needCompression bool, e error) {
	resp := variables.ServResponses{}

	if status == http.StatusOK {

		resp = variables.ServResponses{
			Result: "Data update succesfully",
			Error:  "",
		}
	} else {
		resp = variables.ServResponses{
			Result: "Unsuccesfully",
			Error:  e.Error(),
		}
	}

	strJSON, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("(sendResponceJSON) Marshal error: " + err.Error()))

		variables.PrinterErr(err, "Marshal error:")
		return
	}

	if needCompression {
		w.Header().Set("Content-Encoding", "gzip")
		strJSON, err = compression.Compress(strJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("# (sendResponceJSON) Compress error : " + err.Error()))
			variables.PrinterErr(err, "(sendResponceJSON) Compress error :")
			return
		}
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(strJSON); err != nil {
		variables.PrinterErr(err, ""+"- Send error")
		return
	}
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

		if value, inMap := storage.MGServer.Get(mName); inMap {
			// мы получили байты теперь преобразуем их в float
			byteToFloat, _ := strconv.ParseFloat(string(value), 64)

			if format {
				answer = fmt.Sprintf("%0.3f", byteToFloat)
			} else {

				answer = strconv.FormatFloat(byteToFloat, 'f', -1, 64)
			}

			st = http.StatusOK
		} else {
			err = fmt.Errorf("не найдено имя %s", mName)
			st = http.StatusNotFound //404
			answer = ""
		}

	case "counter":
		if value, inMap := storage.MCServer.Get(mName); inMap {
			// мы получили байты теперь преобразуем их в int
			byteToInt, _ := strconv.ParseInt(string(value), 10, 64)

			if format {
				answer = fmt.Sprintf("%d", byteToInt)
			} else {
				answer = strconv.FormatInt(byteToInt, 10)
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

func readBodyJSONRequest(w http.ResponseWriter, r *http.Request, resp *variables.Metrics, needCompression *bool, doVerificationHash bool) (int, error) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		*needCompression = true
	}

	body, err := io.ReadAll(r.Body)
	variables.PrinterErr(err, "(readBodyJSONRequest) Ошибка чтения тела запроса : ")

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if *needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#(readBodyJSONRequest)ошибка декомпрессии: ")
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	variables.FShowLog(fmt.Sprintf(" Сервер получил json %s \n выполняет дальнейшую обработку \n", string(body)))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		variables.PrinterErr(err, "can't unmarshal: ")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return http.StatusInternalServerError, err
	}

	// Проверим пришла ли хэш функция и если она пришла но не равна получаемой на сервере
	// то отправляем лесом

	if doVerificationHash {
		err = verificationHash(*resp, internal.ConfS.Key)
	}

	if err != nil {
		variables.PrinterErr(err, "")
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil

}

func readBodySliceJSONRequest(w http.ResponseWriter, r *http.Request, resp *[]variables.Metrics, needCompression *bool) (int, error) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		*needCompression = true
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		variables.PrinterErr(err, "#(readBodySliceJSONRequest) ошибка чтения тела запроса: ")

		return http.StatusInternalServerError, err
	}

	if *needCompression {
		body, err = compression.Decompress(body)
		variables.PrinterErr(err, "#HandleGetMetricJSON mistake decompression: ")

	}

	variables.FShowLog(fmt.Sprintf(" Сервер получил json %s \n выполняет дальнейшую обработку \n", string(body)))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		variables.PrinterErr(err, "can't unmarshal: ")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
