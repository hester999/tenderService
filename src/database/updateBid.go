package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
)

func UpdateBid(db *sql.DB, bid *models.Bid, user string) (uuid.UUID, error) {
	current := models.Bid{}

	if err := getBidData(db, bid.Id, &current); err != nil {
		return uuid.Nil, fmt.Errorf("failed to get bid data: %v", err)
	}

	if err := saveCurrentBidToHistory(db, &current); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save current bid to history: %v", err)
	}

	if bid.Name == "" {
		bid.Name = current.Name
	}
	if bid.Description == "" {
		bid.Description = current.Description
	}
	if bid.Status == "" {
		bid.Status = current.Status
	}
	if bid.AuthorType == "" {
		bid.AuthorType = current.AuthorType
	}

	bid.CreatorUserId = current.CreatorUserId
	bid.TenderId = current.TenderId
	bid.AuthorId = current.AuthorId

	current.Version++
	bid.Version = current.Version
	bid.UpdatedAt = time.Now()

	query := `UPDATE bid SET name=$1, description=$2, status=$3, author_type=$4, version=$5, updated_at=$6 WHERE id=$7 RETURNING id`
	var bidID uuid.UUID
	if err := db.QueryRow(query, bid.Name, bid.Description, bid.Status, bid.AuthorType, bid.Version, bid.UpdatedAt, bid.Id).Scan(&bidID); err != nil {
		return uuid.Nil, fmt.Errorf("failed to update bid: %v", err)
	}

	return bidID, nil
}
