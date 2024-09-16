package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"src/database"
	"src/internal"
	"src/models"
	"time"
)

func BidRollBack(w http.ResponseWriter, r *http.Request, db *sql.DB, version int) {
	w.Header().Set("Content-Type", "application/json")
	var bid models.Bid
	strID := internal.GetIdFromURL(r.URL.Path)
	id, _ := uuid.Parse(strID)
	user := r.URL.Query().Get("username")
	err := ValidateBid(db, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "tender not found"})
		return
	}

	err = ValidateUser(db, user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "user not found"})
		return
	}

	err = database.BidRollBack(db, version, id, &bid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}
	response := models.BidResponse{
		Id:         bid.Id.String(),
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId.String(),
		Version:    bid.Version,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
