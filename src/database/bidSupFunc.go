package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
)

func getBidData(db *sql.DB, id uuid.UUID, bid *models.Bid) error {
	query := `SELECT id, tender_id, creator_user_id, name, description, status, author_type, author_id, version, created_at 
              FROM bid 
              WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&bid.Id,
		&bid.TenderId,
		&bid.CreatorUserId,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error fetching bid data: %v", err)
	}

	return nil
}

func saveCurrentBidToHistory(db *sql.DB, bid *models.Bid) error {

	query := `INSERT INTO bid_history 
              (bid_id, tender_id, creator_user_id, name, description, status, author_type, author_id, version, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(query, bid.Id, bid.TenderId, bid.CreatorUserId, bid.Name, bid.Description,
		bid.Status, bid.AuthorType, bid.AuthorId, bid.Version, bid.CreatedAt, time.Now())

	if err != nil {
		return fmt.Errorf("error saving bid history: %v", err)
	}

	return nil
}

func GetBidDataFromHistory(db *sql.DB, id uuid.UUID, bid *models.Bid, version int) error {
	query := `SELECT bid_id, tender_id, creator_user_id, name, description, status, author_type, author_id, version, created_at 
              FROM bid_history 
              WHERE bid_id = $1 and version = $2`

	err := db.QueryRow(query, id, version).Scan(
		&bid.Id,
		&bid.TenderId,
		&bid.CreatorUserId,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error fetching bid history data: %v", err)
	}

	return nil
}
