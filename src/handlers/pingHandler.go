package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"src/models"
)

type Ping struct {
	Message string `json:"message"`
}

func PingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	err := db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Failed to connect to the database"})
		return
	}

	response := Ping{
		Message: "ok",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
