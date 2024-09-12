package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/internal"
	"src/models"
	"time"
)

func AddNewTender(db *sql.DB, tender *models.Tender) (uuid.UUID, error) {

	err := internal.ValidatePermission(db, tender.CreatorUsername, tender.OrganizationId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("permission denied: %v", err)
	}

	err = internal.ValidateUser(db, tender.CreatorUsername)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user does not exist: %v", err)
	}

	var creatorUserID uuid.UUID
	err = db.QueryRow("SELECT id FROM employee WHERE username = $1", tender.CreatorUsername).Scan(&creatorUserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to find user with username %s: %v", tender.CreatorUsername, err)
	}

	if tender.Version == 0 {
		tender.Version = 1
	}
	query := `
        INSERT INTO tender 
        (name, description, service_type, status, version, organization_id, creator_user_id, creator_username) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id,created_at
    `
	var tenderID uuid.UUID
	var createdAt time.Time
	err = db.QueryRow(query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.Version, tender.OrganizationId, creatorUserID, tender.CreatorUsername).Scan(&tenderID, &createdAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert tender: %v", err)
	}

	tender.CreatedAt = createdAt.Format(time.RFC3339)

	return tenderID, nil
}
