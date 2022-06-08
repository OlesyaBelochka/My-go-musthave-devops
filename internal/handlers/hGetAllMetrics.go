package handlers

import (
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/compression"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"log"
	"net/http"
	"strings"
)

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

		data, err := compression.Compress([]byte(html))
		if err != nil {
			variables.PrinterErr(err, "(HandleGetAllMetrics) Ошибка сжатия : ")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		if _, err := w.Write(data); err != nil {
			variables.PrinterErr(err, "(HandleGetAllMetrics) Ошибка отправки : ")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := w.Write([]byte(html)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			variables.PrinterErr(err, "(HandleGetAllMetrics) Ошибка отправки : ")
			return
		}
	}

	w.WriteHeader(http.StatusOK)

}
