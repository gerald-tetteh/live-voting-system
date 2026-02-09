package db

import (
	"database/sql"
	"fmt"
	"os"

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

	createTables(DB, logger)
	return DB
}

func GetDBUrl(logger *zap.Logger) string {
	password, err := os.ReadFile(os.Getenv("DB_PASSWORD_FILE"))
	if err != nil {
		logger.Error("Could not parse credentials")
		logger.Fatal(err.Error())
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", 
		os.Getenv("DB_USER"), 
		password, 
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_NAME"), 
		os.Getenv("DB_SSL_MODE"),
	)
}

func createTables(DB *sql.DB, logger *zap.Logger) {
	createSchema := `
	CREATE TABLE IF NOT EXISTS elections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'closed', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		first_name VARCHAR(20) NOT NULL,
		last_name VARCHAR(20) NOT NULL,
		middle_name VARCHAR(20),
		email VARCHAR(50) UNIQUE NOT NULL,
		role VARCHAR(20) DEFAULT 'base' CHECK (role IN ('base', 'admin')),
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS votes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		election_id UUID REFERENCES elections(id),
		user_id UUID REFERENCES users(id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(createSchema);

	if err != nil {
		logger.Error("Could not create tables")
		logger.Fatal(err.Error())
	}
}