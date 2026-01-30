package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func InitDB(logger *zap.Logger, dbUrl string, maxOpenConnections, maxIdleConnections int) *sql.DB {
	DB, err := sql.Open("postgres", dbUrl)
	if err != nil {
		logger.Error("Could not connect to database")
		logger.Fatal(err.Error())
	}

	DB.SetMaxOpenConns(maxOpenConnections)
	DB.SetMaxIdleConns(maxIdleConnections)
	return DB
}