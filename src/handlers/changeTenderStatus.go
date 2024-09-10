package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"src/models"
	"strconv"
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

func getTenderId(url string) (int, error) {

	url = strings.TrimPrefix(url, "/api/tenders/")
	url = strings.TrimSuffix(url, "/status")

	id, err := strconv.Atoi(url)
	if err != nil || id < 1 {
		return http.StatusBadRequest, fmt.Errorf("invalid service id provided") // возможно будут  баги
	}
	return id, nil

}
func ChangeTenderStatus(r http.Request, w http.ResponseWriter, db *sql.DB) {

	//response := models.TenderResponse{}
	w.Header().Set("Content-Type", "application/json")

	avaliibleTenderStatus := map[string]bool{
		"Created":   true,
		"Published": true,
		"Closed":    true,
	}

	id, err := getTenderId(r.URL.Path)

	if err != nil {
		w.WriteHeader(id)
		errResp := models.ErrorResponse{Reason: err.Error()}
		json.NewEncoder(w).Encode(errResp)
		return
	}
	status := r.URL.Query().Get("status")
	user := r.URL.Query().Get("username")
	err = validateStatus([]string{status}, avaliibleTenderStatus)
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

	//err := ChangeTenderStatus(db, &response, status, user)
}
