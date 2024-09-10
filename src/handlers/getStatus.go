package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"src/database"
	"src/models"
	"strings"
)

func GetBannerStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Извлекаем ID тендера из URL
	strId := r.URL.Path
	strId = strings.TrimPrefix(strId, "/api/tenders/")
	strId = strings.TrimSuffix(strId, "/status")

	// Извлекаем username из параметров запроса
	user := r.URL.Query().Get("username")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "username is required"})
		return
	}

	// Преобразуем ID тендера в UUID
	tenderId, err := uuid.Parse(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid tender ID format"})
		return
	}

	// Получаем статус тендера
	status, err := database.GetTenderStatus(db, tenderId, user)
	if err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "user not found"})
			return
		} else if err.Error() == "tender not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "tender not found"})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "internal server error"})
			return
		}
	}

	// Возвращаем статус тендера в ответе

	json.NewEncoder(w).Encode(status)
}
