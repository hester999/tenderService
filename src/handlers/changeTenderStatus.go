package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"src/database"
	"src/models"
	"strings"
)

func validateStatus(tenderStatus []string, validTenderStatus map[string]bool) error {
	for _, tenderStatus := range tenderStatus {
		if _, exists := validTenderStatus[tenderStatus]; !exists {
			return fmt.Errorf("Invalid tender status  provided. Available values: Created, Published, Closed")
		}
	}
	return nil
}

func getTenderId(url string) string {

	url = strings.TrimPrefix(url, "/api/tenders/")
	url = strings.TrimSuffix(url, "/status")

	return url

}
func ChangeTenderStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	resp := models.TenderResponse{}
	w.Header().Set("Content-Type", "application/json")

	avaliibleTenderStatus := map[string]bool{
		"Created":   true,
		"Published": true,
		"Closed":    true,
	}

	idStr := getTenderId(r.URL.Path)
	id, _ := uuid.Parse(idStr)

	status := r.URL.Query().Get("status")
	user := r.URL.Query().Get("username")
	err := validateStatus([]string{status}, avaliibleTenderStatus)
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
