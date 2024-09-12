package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

func GetTenderStatus(db *sql.DB, tenderId uuid.UUID, creator string) (string, error) {

	var userExists bool
	userQuery := "SELECT EXISTS(SELECT 1 FROM employee WHERE username = $1)"
	err := db.QueryRow(userQuery, creator).Scan(&userExists)
	if err != nil {
		return "", fmt.Errorf("internal error: failed to check user: %v", err)
	}
	if !userExists {
		return "", fmt.Errorf("user not found")
	}

	query := "SELECT status FROM tender WHERE id = $1 AND creator_username = $2"
	row := db.QueryRow(query, tenderId, creator)

	var status string
	if err := row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("tender not found")
		}
		return "", fmt.Errorf("failed to retrieve tender status: %v", err)
	}

	return status, nil
}
