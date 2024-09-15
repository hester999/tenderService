package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
)

func checkTender(db *sql.DB, id uuid.UUID) bool {
	query := "select exists(select id from tender where id = $1)"

	var status bool
	_ = db.QueryRow(query, id).Scan(&status)
	if status == false {

		return status
	}
	return true
}

func getTenderData(db *sql.DB, id uuid.UUID, tender *models.Tender) error {

	query := `SELECT id, name, description, service_type, status, organization_id, creator_user_id, creator_username, version, created_at 
              FROM tender 
              WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUserId,
		&tender.CreatorUsername,
		&tender.Version,
		&tender.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error fetching tender data: %v", err)
	}

	return nil
}

func saveCurrentTenderToHistory(db *sql.DB, tender *models.Tender) error {

	query := `INSERT INTO tender_history 
              (tender_id, name, description, service_type, status, version, organization_id, creator_user_id, creator_username, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(query, tender.Id, tender.Name, tender.Description, tender.ServiceType, tender.Status,
		tender.Version, tender.OrganizationId, tender.CreatorUserId, tender.CreatorUsername, tender.CreatedAt, time.Now())
	return err
}

func GetDataFromHistoryTender(db *sql.DB, id uuid.UUID, tender *models.Tender, version int) error {
	query := `SELECT tender_id, name, description, service_type, status, organization_id, creator_user_id, creator_username, version, created_at 
              FROM tender_history 
              WHERE tender_id = $1 and version = $2`

	err := db.QueryRow(query, id, version).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUserId,
		&tender.CreatorUsername,
		&tender.Version,
		&tender.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error fetching tender data: %v", err)
	}

	return nil
}
