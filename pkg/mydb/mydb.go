package mydb

import (
	"database/sql"

	_ "github.com/lib/pq" // Импортируем PostgreSQL драйвер
)

func Init(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
