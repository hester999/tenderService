package internal

import (
	"database/sql"
	"errors"
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

func ValidateUserBid(db *sql.DB, userID uuid.UUID) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT id FROM employee WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if !exists {
		return errors.New("user not found")
	}

	var count int
	err = db.QueryRow("SELECT COUNT(user_id) FROM organization_responsible WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking user permissions: %w", err)
	}

	if count == 0 {
		return errors.New("user has no responsibilities")
	}

	return nil
}

func ValidateTenderId(db *sql.DB, tenderId uuid.UUID) error {
	var exist bool
	err := db.QueryRow(`SELECT EXISTS(SELECT id FROM tender WHERE id = $1)`, tenderId).Scan(&exist)

	if err != nil {
		return fmt.Errorf("error checking tender existence: %w", err)
	}

	if !exist {
		return errors.New("teder not found")
	}

	return nil

}
