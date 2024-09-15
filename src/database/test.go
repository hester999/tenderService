package database

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Employee struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type Organization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type OrganizationResponsible struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
}

func GetAllEmployees(db *sql.DB) ([]Employee, error) {
	query := `SELECT id, username, first_name, last_name, created_at, updated_at FROM employee`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Username, &e.FirstName, &e.LastName, &e.CreatedAt, &e.UpdatedAt); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func GetAllOrganizations(db *sql.DB) ([]Organization, error) {
	query := `SELECT id, name, description, type, created_at, updated_at FROM organization`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var organizations []Organization
	for rows.Next() {
		var o Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.Description, &o.Type, &o.CreatedAt, &o.UpdatedAt); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		organizations = append(organizations, o)
	}
	return organizations, nil
}

func GetAllOrganizationResponsibles(db *sql.DB) ([]OrganizationResponsible, error) {
	query := `SELECT id, organization_id, user_id FROM organization_responsible`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var responsibles []OrganizationResponsible
	for rows.Next() {
		var r OrganizationResponsible
		if err := rows.Scan(&r.ID, &r.OrganizationID, &r.UserID); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		responsibles = append(responsibles, r)
	}
	return responsibles, nil
}

func SendAllData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	employees, err := GetAllEmployees(db)
	if err != nil {
		http.Error(w, "Failed to get employees", http.StatusInternalServerError)
		return
	}

	organizations, err := GetAllOrganizations(db)
	if err != nil {
		http.Error(w, "Failed to get organizations", http.StatusInternalServerError)
		return
	}

	responsibles, err := GetAllOrganizationResponsibles(db)
	if err != nil {
		http.Error(w, "Failed to get organization responsibles", http.StatusInternalServerError)
		return
	}

	response := struct {
		Employees                []Employee                `json:"employees"`
		Organizations            []Organization            `json:"organizations"`
		OrganizationResponsibles []OrganizationResponsible `json:"organization_responsibles"`
	}{
		Employees:                employees,
		Organizations:            organizations,
		OrganizationResponsibles: responsibles,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
