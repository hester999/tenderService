package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
)

func CreateNewBid(bid *models.Bid, db *sql.DB) error {

	if bid.Version == 0 {
		bid.Version = 1
	}

	query := `
		INSERT INTO bid (
			tender_id,
			creator_user_id,
			name,
			description,
			status,
			author_type,  
			author_id,    
			version
		) 
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) 
		RETURNING id, created_at;
	`

	var newBidID uuid.UUID
	var createdAt time.Time

	err := db.QueryRow(
		query,
		bid.TenderId,
		bid.CreatorUserId,
		bid.Name,
		bid.Description,
		bid.Status,
		bid.AuthorType,
		bid.AuthorId,
		bid.Version,
	).Scan(&newBidID, &createdAt)

	if err != nil {
		return fmt.Errorf("failed to create new bid: %v", err)
	}

	bid.Id = newBidID
	bid.CreatedAt = createdAt

	return nil
}
