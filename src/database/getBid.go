package database

import (
	"database/sql"
	"fmt"
	"src/models"
)

func GetBid(bid *models.ResponseGetBid, db *sql.DB) error {
	query := `
		SELECT id, tender_id, creator_user_id, name, description, status, author_type, author_id, version, created_at, updated_at
		FROM bid
		WHERE creator_user_id = (SELECT id FROM employee WHERE username = $1)
		LIMIT $2 OFFSET $3;
	`

	rows, err := db.Query(query, bid.User, bid.Limit, bid.Offset)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Bid
		if err := rows.Scan(
			&b.Id, &b.TenderId, &b.CreatorUserId, &b.Name, &b.Description, &b.Status, &b.AuthorType, &b.AuthorId, &b.Version, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		bid.Bid = append(bid.Bid, b)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}

	return nil
}
