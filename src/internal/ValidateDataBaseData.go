package internal

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

func ValidatePermission(db *sql.DB, username string, organizationID uuid.UUID) error {
	var responsibleID uuid.UUID
	query := `
        SELECT id 
        FROM organization_responsible 
        WHERE user_id = (SELECT id FROM employee WHERE username = $1) 
        AND organization_id = $2
    `
	err := db.QueryRow(query, username, organizationID).Scan(&responsibleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %s is not responsible for organization %s", username, organizationID)
		}
		return fmt.Errorf("failed to validate permission: %v", err)
	}

	return nil
}

func ValidateUser(db *sql.DB, username string) error {
	query := `SELECT username FROM employee WHERE username = $1`

	err := db.QueryRow(query, username).Scan(&username)

	if err == sql.ErrNoRows {
		return fmt.Errorf("user does not exist")
	}
	return nil
}
