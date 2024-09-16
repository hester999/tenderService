package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"src/database"

	"github.com/google/uuid"
	"net/http"
	"src/internal"
	"src/models"
)

func ValidateUser(db *sql.DB, user string) error {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM employee WHERE username=$1)`
	err := db.QueryRow(query, user).Scan(&exist)

	if err != nil {
		return err
	}
	if !exist {
		return errors.New("User does not exist")
	}

	query = `SELECT COUNT(user_id)
	FROM organization_responsible
	WHERE user_id = (SELECT id FROM employee WHERE username = $1)`

	var count int
	_ = db.QueryRow(query, user).Scan(&count)

	if count != 1 {
		return errors.New("user must be responsible for only one organization")
	}

	return nil
}

func ValidateBid(db *sql.DB, id uuid.UUID) error {
	var exist bool
	query := `SELECT exists(SELECT 1 FROM bid WHERE id=$1)`

	err := db.QueryRow(query, id).Scan(&exist)

	if err != nil {
		return err
	}

	if !exist {
		return errors.New("bid not found")
	}
	return nil
}

func ChangeBidStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	resp := models.BidResponse{}

	w.Header().Set("Content-Type", "application/json")

	avaliibleBidStatus := map[string]bool{
		"Created":   true,
		"Published": true,
		"Closed":    true,
	}

	idStr := internal.GetIdFromURL(r.URL.Path)
	id, _ := uuid.Parse(idStr)

	status := r.URL.Query().Get("status")
	user := r.URL.Query().Get("username")
	if user == "" || status == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "status or user can not be empty"})
		return
	}
	err := ValidateUser(db, user)
	if err != nil {
		switch err.Error() {
		case "User does not exist":
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		case "user must be responsible for only one organization":
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "user must be responsible for only one organization" +
				""})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
			return
		}
	}
	err = ValidateBid(db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "bid not found"})
		return
	}

	err = internal.ValidateStatus([]string{status}, avaliibleBidStatus)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid bid status  provided. Available values: Created, Published, Closed"})
		return
	}

	err = database.ChangeBidStatus(db, id, status, &resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(resp)
}
