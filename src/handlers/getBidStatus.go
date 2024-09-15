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

func GetBidStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	strId := internal.GetTenderId(r.URL.Path)

	user := r.URL.Query().Get("username")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "username is required"})
		return
	}

	BidID, err := uuid.Parse(strId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid tender ID format"})
		return
	}

	status, err := database.GetBidData(BidID, db)

	if err != nil {
		switch err.Error() {
		case "the user does not have enough rights":
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "the user does not have enough rights"})
			return
		case "not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "not found"})
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "internal server error"})
			return
		}
	}
	json.NewEncoder(w).Encode(status)

}
