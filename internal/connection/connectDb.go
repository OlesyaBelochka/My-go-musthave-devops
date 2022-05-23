package connection

import (
	"database/sql"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal"
	_ "github.com/jackc/pgx/stdlib"
)

func Start(config *internal.ConfigServer) error {

	db, err := sql.Open("pgx", config.DatabaseURL)

	if err != nil {
		return err
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
