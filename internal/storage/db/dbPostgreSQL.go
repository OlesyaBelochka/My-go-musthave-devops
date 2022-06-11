package db

import (
	"context"
	"database/sql"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/config"
	_ "github.com/jackc/pgx/stdlib"
)

const (
	// TODO(spencer): update the CREATE DATABASE statement in the schema
	//   to pull out the database specified in the DB URL and use it instead
	//   of "photos" below.
	photosSchema = `
CREATE TABLE IF NOT EXISTS metrics (
 id           TEXT,
 mtype 	  TEXT,
 delta		   BIGINT,
 val        DOUBLE PRECISION,
 PRIMARY KEY (id, mtype)
);`
)

func OpenDB(config *config.ConfigServer) (*sql.DB, error) {
	return sql.Open("pgx", config.DatabaseURL)
}

func InitSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, photosSchema)
	return err
}
