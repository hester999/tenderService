package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"src/database"
	"src/models"
)

func ChangeTender(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	idStr := getTenderId(r.URL.Path)
	id, _ := uuid.Parse(idStr)

	tender := models.Tender{}
	resp := models.TenderResponse{}

	err := json.NewDecoder(r.Body).Decode(&tender)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid request payload"})
		return
	}

	err = tender.ValidateChangeFiled()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	err = database.ChangeTender(db, id, &tender, &resp)
	if err != nil {
		switch err.Error() {
		case "tender not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "tender not found"})
		}

	}
	return
}
