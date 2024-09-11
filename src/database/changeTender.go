package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
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

	// Выполнение запроса и сканирование результатов непосредственно в структуру tender
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

func ChangeTender(db *sql.DB, id uuid.UUID, tender *models.Tender, user string) (uuid.UUID, error) {
	tempTender := models.Tender{}
	if !checkTender(db, id) {
		return uuid.Nil, fmt.Errorf("tender not found")
	}

	if err := getTenderData(db, id, &tempTender); err != nil {
		return uuid.Nil, fmt.Errorf("failed to get tender data: %v", err)
	}

	if tender.Name == "" {
		tender.Name = tempTender.Name
	}
	if tender.Description == "" {
		tender.Description = tempTender.Description
	}
	if tender.Status == "" {
		tender.Status = tempTender.Status
	}
	if tender.ServiceType == "" {
		tender.ServiceType = tempTender.ServiceType
	}

	tender.CreatorUsername = user

	tender.OrganizationId = tempTender.OrganizationId
	tender.CreatorUserId = tempTender.CreatorUserId

	tender.Version = tempTender.Version + 1
	tenderID, err := AddNewTender(db, tender)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add new tender version: %v", err)
	}

	return tenderID, nil
}
