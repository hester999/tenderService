package database

import (
	"database/sql"
	"errors"
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

func GetBidData(id uuid.UUID, db *sql.DB) (string, error) {
	var status string
	var userPermission bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE user_id = (SELECT author_id FROM bid WHERE id = $1))`, id).Scan(&userPermission)
	if !userPermission {
		return "", errors.New("the user does not have enough rights")
	}

	err = db.QueryRow("SELECT status FROM bid WHERE id = $1", id).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found")
		}
	}

	return status, nil
}
