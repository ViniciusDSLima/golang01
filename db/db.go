package db

import (
	"database/sql"
	"fmt"
	"github.com/ViniciusDSLima/golang01/config"
)

func NewPostgresStorage(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
