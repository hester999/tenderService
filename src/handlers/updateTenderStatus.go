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

func ChangeTenderStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	resp := models.TenderResponse{}
	w.Header().Set("Content-Type", "application/json")

	avaliibleTenderStatus := map[string]bool{
		"Created":   true,
		"Published": true,
		"Closed":    true,
	}

	idStr := internal.GetIdFromURL(r.URL.Path)
	id, _ := uuid.Parse(idStr)

	status := r.URL.Query().Get("status")
	user := r.URL.Query().Get("username")

	err := internal.ValidateStatus([]string{status}, avaliibleTenderStatus)
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		errResp := models.ErrorResponse{Reason: "Invalid username provided"}
		json.NewEncoder(w).Encode(errResp)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResp := models.ErrorResponse{Reason: err.Error()}
		json.NewEncoder(w).Encode(errResp)
		return
	}

	//(db *sql.DB, response *models.TenderResponse, status, user string
	err = database.ChangeTenderStatus(db, &resp, id, user, status)
	if err != nil {
		switch err.Error() {
		case "permission denied":
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Insufficient permissions to perform this action."})

		case "user not found or not exist":
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "User not found or incorrect."})

		case "tender not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Tender not found."})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Internal server error."})
		}
		return
	}

	json.NewEncoder(w).Encode(resp)

}
