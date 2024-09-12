package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"src/database"
	"src/internal"
	"src/models"
)

// todo доделать ообработку всех ошибок
func TenderRollBack(w http.ResponseWriter, r *http.Request, db *sql.DB, version int) {
	tender := models.Tender{}
	w.Header().Set("Content-Type", "application/json")

	idStr := internal.GetTenderId(r.URL.Path)
	id, _ := uuid.Parse(idStr)
	fmt.Println(id)
	user := r.URL.Query().Get("username")

	tenderVersionFromDateBase := internal.GetVersionFromHistoryTender(db, id, version)

	if !tenderVersionFromDateBase {
		fmt.Println(version, tenderVersionFromDateBase)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "version no found"})
		return
	}

	err := internal.ValidateUser(db, user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: "user no found"})
		return
	}

	err = database.TenderRollBack(db, id, user, &tender, version)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Reason: err.Error()})
		return
	}
	resp := models.TenderResponse{
		Id:          tender.Id.String(),
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      tender.Status,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}

	json.NewEncoder(w).Encode(resp)

}
