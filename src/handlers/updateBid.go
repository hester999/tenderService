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

func UpdateBid(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	w.Header().Set("Content-Type", "application/json")

	var bid models.Bid

	// Извлекаем ID из URL
	strID := internal.GetIdFromURL(r.URL.Path)
	user := r.URL.Query().Get("username")

	if strID == "" || user == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "user or bidID is empty"})
		return
	}

	err := json.NewDecoder(r.Body).Decode(&bid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid request payload"})
		return
	}

	bid.Id, err = uuid.Parse(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid bid ID"})
		return
	}

	err = bid.ValidateChangeFields()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	err = internal.ValidateUser(db, user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	err = ValidateBid(db, bid.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "bid not found"})
		return
	}

	updatedBidID, err := database.UpdateBid(db, &bid, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "internal service error"})
		return
	}

	response := models.BidResponse{
		Id:         updatedBidID.String(),
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
