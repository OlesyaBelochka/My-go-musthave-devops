package handlers

import (
	"context"
	config "github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/storage/db"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
	"net/http"
	"time"
)

func HandlePingDB(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	dataBase, err := db.OpenDB(config.ConfS)

	if err != nil {
		variables.PrinterErr(err, "#(HandlePingDB) ошибка при открытии БД: ")
		sendStatus(w, http.StatusInternalServerError)
	}

	defer func() { _ = dataBase.Close() }()

	if err := db.InitSchema(ctx, dataBase); err != nil {
		variables.PrinterErr(err, "#(HandlePingDB)ошибка при создании схемы инициализации: ")
		sendStatus(w, http.StatusInternalServerError)
	}

	sendStatus(w, http.StatusOK)

}
