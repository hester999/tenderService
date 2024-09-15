package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"src/database"
	"src/models"
	"strconv"
)

func GetBid(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	var bid models.ResponseGetBid
	bid.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	bid.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	bid.User = r.URL.Query().Get("username")

	err := bid.ValidateUser(bid.User, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	err = bid.ValidatePermission(bid.User, db)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	if bid.Offset < 0 || bid.Limit < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid pagination parameters"})
		return
	}

	err = database.GetBid(&bid, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	var bidResponses []models.BidResponse
	for _, b := range bid.Bid {
		bidResponses = append(bidResponses, models.BidResponse{
			Id:         b.Id.String(),
			Name:       b.Name,
			Status:     b.Status,
			AuthorType: b.AuthorType,
			AuthorId:   b.AuthorId.String(),
			Version:    b.Version,
			CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	json.NewEncoder(w).Encode(bidResponses)
}
