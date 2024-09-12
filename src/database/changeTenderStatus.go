package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/internal"
	"src/models"
	"time"
)

func ChangeTenderStatus(db *sql.DB, tender *models.TenderResponse, id uuid.UUID, user, status string) error {
	var organizationID uuid.UUID

	err := db.QueryRow(`SELECT organization_id FROM tender WHERE id = $1`, id).Scan(&organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("tender not found")
		}
		return fmt.Errorf("failed to fetch organization id: %v", err)
	}

	err = internal.ValidatePermission(db, user, organizationID)
	if err != nil {
		return fmt.Errorf("permission denied: %v", err)
	}

	err = internal.ValidateUser(db, user)
	if err != nil {
		return fmt.Errorf("user not found or does not exist")
	}

	query := `
		UPDATE tender 
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, name, description, service_type, status, version, updated_at
	`

	var updatedAt time.Time
	err = db.QueryRow(query, status, id).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("tender not found")
		}
		return fmt.Errorf("failed to update tender status: %v", err)
	}

	tender.CreatedAt = updatedAt.Format(time.RFC3339)

	return nil
}
