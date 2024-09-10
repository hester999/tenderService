package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"src/database"
	"src/models"
	"strconv"
	_ "strings"
)

// добавить проверку на все аргументы в url
func validateServiceTypes(serviceTypes []string, validServiceTypes map[string]bool) error {
	for _, serviceType := range serviceTypes {
		if _, exists := validServiceTypes[serviceType]; !exists {
			return fmt.Errorf("Invalid service type provided. Available values: Construction, Delivery, Manufacture")
		}
	}
	return nil
}

func GetTender(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	validServiceTypes := map[string]bool{
		"Construction": true,
		"Delivery":     true,
		"Manufacture":  true,
	}

	w.Header().Set("Content-Type", "application/json")

	tender := models.ResponseGetBanner{}
	tender.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	tender.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))

	username := r.URL.Query().Get("username")
	serviceTypes := r.URL.Query()["service_type"]

	tender.Validatepaging()

	if tender.Offset < 0 || tender.Limit < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Invalid pagination parameters"})
		return
	}

	var err error
	if username != "" {
		tender.User = username
		err = database.GetTenderByUser(db, &tender)
	} else if len(serviceTypes) > 0 {
		if err := validateServiceTypes(serviceTypes, validServiceTypes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
			return
		}
		err = database.GetTenderByServiceType(db, &tender, serviceTypes)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Either username or service_type must be provided"})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "Failed to get tender"})
		return
	}

	var response []models.TenderResponse
	for _, t := range tender.Tenders {
		// UUID уже имеет правильный тип, поэтому конвертация в строку происходит здесь
		tr := models.TenderResponse{
			Id:          t.Id.String(), // Преобразуем UUID в строку
			Name:        t.Name,
			Description: t.Description,
			Status:      t.Status,
			ServiceType: t.ServiceType,
			Version:     t.Version,
			CreatedAt:   t.CreatedAt,
		}
		response = append(response, tr)
	}

	json.NewEncoder(w).Encode(response)
}
