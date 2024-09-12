package database

import (
	"database/sql"
	"fmt"
	"src/models"
	"strings"
)

func GetTenderByUser(db *sql.DB, tender *models.ResponseGetBanner) error {
	query := `
		SELECT id, name, description, status, service_type, version, created_at 
		FROM tender 
		WHERE creator_username = $1 
		ORDER BY name 
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, tender.User, tender.Limit, tender.Offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	tender.Tenders = []models.Tender{}

	for rows.Next() {
		var t models.Tender

		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.Version, &t.CreatedAt)
		if err != nil {
			return err
		}

		tender.Tenders = append(tender.Tenders, t)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func GetTenderByServiceType(db *sql.DB, tender *models.ResponseGetBanner, serviceTypes []string) error {

	placeholders := make([]string, len(serviceTypes))
	for i := range serviceTypes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, status, service_type, version, created_at 
		FROM tender 
		WHERE service_type IN (%s)
		ORDER BY name 
		LIMIT $%d OFFSET $%d
	`, strings.Join(placeholders, ", "), len(serviceTypes)+1, len(serviceTypes)+2)

	args := make([]interface{}, len(serviceTypes)+2)
	for i, v := range serviceTypes {
		args[i] = v
	}
	args[len(serviceTypes)] = tender.Limit
	args[len(serviceTypes)+1] = tender.Offset

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	tender.Tenders = []models.Tender{}

	for rows.Next() {
		var t models.Tender
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.Version, &t.CreatedAt)
		if err != nil {
			return err
		}
		tender.Tenders = append(tender.Tenders, t)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
