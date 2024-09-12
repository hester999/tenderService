package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host    = "localhost"
	port    = 5432
	user    = "admin"
	pass    = "123"
	db_name = "avito"
)

func InitDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, db_name)

	db, _ := sql.Open("postgres", psqlInfo)

	return db, nil
}

//package database
//
//import (
//	"database/sql"
//	"fmt"
//	"os"
//
//	_ "github.com/lib/pq"
//)
//
//func InitDB() (*sql.DB, error) {
//
//	host := os.Getenv("POSTGRES_HOST")
//	port := os.Getenv("POSTGRES_PORT")
//	user := os.Getenv("POSTGRES_USERNAME")
//	pass := os.Getenv("POSTGRES_PASSWORD")
//	dbName := os.Getenv("POSTGRES_DATABASE")
//
//	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
//		host, port, user, pass, dbName)
//
//	db, err := sql.Open("postgres", psqlInfo)
//	if err != nil {
//		return nil, fmt.Errorf("Failed to open a DB connection: %v", err)
//	}
//	err = db.Ping()
//	if err != nil {
//		return nil, fmt.Errorf("Failed to ping DB: %v", err)
//	}
//
//	return db, nil
//}
