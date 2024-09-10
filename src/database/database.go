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
