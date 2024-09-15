package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
)

func ChangeBidStatus(db *sql.DB, id uuid.UUID, status string, resp *models.BidResponse) error {

	query := `
		UPDATE bid
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, name, status, author_type, author_id, version, created_at;
	`

	var updatedAt time.Time

	err := db.QueryRow(query, status, id).Scan(
		&resp.Id,
		&resp.Name,
		&resp.Status,
		&resp.AuthorType,
		&resp.AuthorId,
		&resp.Version,
		&updatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update bid status: %v", err)
	}

	resp.CreatedAt = updatedAt.Format(time.RFC3339)

	return nil
}
