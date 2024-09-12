package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"src/database"
	"src/internal"
	"src/models"
)

func ChangeTender(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	idStr := internal.GetTenderId(r.URL.Path)
	id, _ := uuid.Parse(idStr)

	tender := models.Tender{}
	resp := models.TenderResponse{}
	user := r.URL.Query().Get("username")

	err := internal.ValidateUser(db, user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "user not found"})
		return
	}
	err = json.NewDecoder(r.Body).Decode(&tender)

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

	tenderID, err := database.ChangeTender(db, id, &tender, user)
	if err != nil {
		switch err.Error() {
		case "tender not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "tender not found"})
			return
		case "permission denied":
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "permission denied"})
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "internal service error"})
			return
		}
	}
	resp = models.TenderResponse{
		Id:          tenderID.String(),
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      tender.Status,
		Version:     tender.Version,
		CreatedAt:   tender.UpdatedAt,
	}
	json.NewEncoder(w).Encode(resp)
	return
}
