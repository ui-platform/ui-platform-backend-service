package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func NewPostgresDB(DBHost, DBPort, DBUser, DBName, DBPass, DBSSLMode string) (*PostgresDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		DBHost, DBPort, DBUser, DBName, DBPass, DBSSLMode)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1000)               // Рассчитано для 10000 RPS при среднем времени обработки запроса
	db.SetMaxIdleConns(100)                // Достаточно для обработки обычной нагрузки
	db.SetConnMaxLifetime(5 * time.Minute) // Уменьшение вероятности проблем из-за старения соединений

	return &PostgresDB{
		DB: db,
	}, nil
}
