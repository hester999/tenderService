package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USERNAME")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DATABASE")
	sslMode := os.Getenv("POSTGRES_SSLMODE")

	if sslMode == "" {
		sslMode = "require"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbName, sslMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to open a DB connection: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping DB: %v", err)
	}
	err = createTables(db)
	if err != nil {
		fmt.Errorf("Failed to create tables")
	}
	return db, nil
}
