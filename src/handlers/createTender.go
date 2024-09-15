package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"src/database"
	"src/models"
)

func CreateTender(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	tender := models.Tender{}

	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid request payload"})
		return
	}

	if err := tender.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	tenderID, err := database.AddNewTender(db, &tender)
	if err != nil {
		if err.Error() == "user does not exist" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "User not found or user does not exist"})
			return
		} else if err.Error() == "permission denied" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Permission denied: user is not responsible for the organization"})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Failed to create tender"})
			return
		}
	}

	response := models.TenderResponse{
		Id:          tenderID.String(), // Преобразование UUID в строку
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      tender.Status,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}

	json.NewEncoder(w).Encode(response)
}
