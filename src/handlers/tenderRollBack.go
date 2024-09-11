package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"src/internal"
	"src/models"
)

func TenderRollBack(w http.ResponseWriter, r *http.Request, db *sql.DB, version int) {

	w.Header().Set("Content-Type", "application/json")

	idStr := internal.GetTenderId(r.URL.Path)
	id, _ := uuid.Parse(idStr)
	fmt.Println(id)
	user := r.URL.Query().Get("username")

	tenderVersionFromDateBase := internal.GetVersion(db, id)

	if version != tenderVersionFromDateBase {
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
}
