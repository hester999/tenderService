package handlers

import (
	"database/sql"
	"encoding/json"
	_ "errors"
	"net/http"
	"src/database"
	"src/models"
	"time"
)

func CreateBids(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var bid models.Bid
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	if err := bid.ValidateFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	if err := bid.ValidateBidUser(db); err != nil {
		switch err.Error() {
		case "user not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "User not found"})
			return
		case "organization not found":
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "A user can only be responsible for one organization"})
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
			return
		}
	}

	if err := bid.ValidateTenderId(db); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	err := database.CreateNewBid(&bid, db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	bidResp := models.BidResponse{
		Id:         bid.Id.String(),
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId.String(),
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt.Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(bidResp)
}
