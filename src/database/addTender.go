package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
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

func AddNewTender(db *sql.DB, tender *models.Tender) (uuid.UUID, error) {

	// Проверяем, есть ли у пользователя права на создание тендера
	err := ValidatePermission(db, tender.CreatorUsername, tender.OrganizationId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("permission denied: %v", err)
	}

	// Проверяем, существует ли пользователь
	err = ValidateUser(db, tender.CreatorUsername)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user does not exist: %v", err)
	}

	// Получаем UUID пользователя по его username
	var creatorUserID uuid.UUID
	err = db.QueryRow("SELECT id FROM employee WHERE username = $1", tender.CreatorUsername).Scan(&creatorUserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to find user with username %s: %v", tender.CreatorUsername, err)
	}

	// Вставляем новый тендер и возвращаем его id и created_at
	query := `
        INSERT INTO tender 
        (name, description, service_type, status, version, organization_id, creator_user_id, creator_username) 
        VALUES ($1, $2, $3, $4, 1, $5, $6, $7)
        RETURNING id, created_at
    `
	var tenderID uuid.UUID
	var createdAt time.Time
	err = db.QueryRow(query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, creatorUserID, tender.CreatorUsername).Scan(&tenderID, &createdAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert tender: %v", err)
	}

	// Присваиваем полученные значения в структуру тендера
	tender.CreatedAt = createdAt.Format(time.RFC3339)
	tender.Version = 1

	return tenderID, nil
}
