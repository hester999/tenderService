package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type DropTablesResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func DropTables(w http.ResponseWriter, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	tx, err := db.Begin()
	if err != nil {
		response := DropTablesResponse{
			Message: "Failed to start transaction",
			Status:  "error",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	dropTableQueries := []string{
		`DROP TABLE IF EXISTS bid_history CASCADE;`,
		`DROP TABLE IF EXISTS bid CASCADE;`,
		`DROP TABLE IF EXISTS tender_history CASCADE;`,
		`DROP TABLE IF EXISTS tender CASCADE;`,
	}

	for _, query := range dropTableQueries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			response := DropTablesResponse{
				Message: fmt.Sprintf("Failed to execute query: %v", query),
				Status:  "error",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		response := DropTablesResponse{
			Message: "Failed to commit transaction",
			Status:  "error",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DropTablesResponse{
		Message: "Tables bid_history, bid, tender_history, and tender have been successfully dropped.",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(response)
}
