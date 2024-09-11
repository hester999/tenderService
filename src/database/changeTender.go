package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
	"time"
)

// Получение данных о тендере и проверка возможности изменения
func checkTender(db *sql.DB, id uuid.UUID) bool {
	query := "select exists(select id from tender where id = $1)"

	var status bool
	_ = db.QueryRow(query, id).Scan(&status)
	if status == false {

		return status
	}
	return true
}
func getTenderData(db *sql.DB, id uuid.UUID, tender *models.Tender) error {

	query := `SELECT id, name, description, service_type, status, organization_id, creator_user_id, creator_username, version, created_at 
              FROM tender 
              WHERE id = $1`

	// Выполнение запроса и сканирование результатов непосредственно в структуру tender
	err := db.QueryRow(query, id).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUserId,
		&tender.CreatorUsername,
		&tender.Version,
		&tender.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error fetching tender data: %v", err)
	}

	return nil
}

func saveCurrentTenderToHistory(db *sql.DB, tender *models.Tender) error {

	query := `INSERT INTO tender_history 
              (tender_id, name, description, service_type, status, version, organization_id, creator_user_id, creator_username, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(query, tender.Id, tender.Name, tender.Description, tender.ServiceType, tender.Status,
		tender.Version, tender.OrganizationId, tender.CreatorUserId, tender.CreatorUsername, tender.CreatedAt, time.Now())
	return err
}

func ChangeTender(db *sql.DB, id uuid.UUID, updates *models.Tender, user string) (uuid.UUID, error) {
	current := models.Tender{}
	if !checkTender(db, id) {
		return uuid.Nil, fmt.Errorf("tender not found")
	}
	if err := getTenderData(db, id, &current); err != nil {
		return uuid.Nil, fmt.Errorf("failed to get tender data: %v", err)
	}

	if err := updates.ValidateChangeFiled(); err != nil {
		return uuid.Nil, err
	}

	// Save current state to history
	if err := saveCurrentTenderToHistory(db, &current); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save current tender to history: %v", err)
	}

	// Update fields only if provided in the request
	if updates.Name == "" {
		updates.Name = current.Name
	}
	if updates.Description == "" {
		updates.Description = current.Description
	}
	if updates.Status == "" {
		updates.Status = current.Status
	}
	if updates.ServiceType == "" {
		updates.ServiceType = current.ServiceType
	}

	updates.CreatorUsername = user

	updates.OrganizationId = current.OrganizationId
	updates.CreatorUserId = current.CreatorUserId
	// Increment version and update updated_at
	current.Version++
	updates.Version = current.Version
	updates.UpdatedAt = time.Now().Format(time.RFC3339)

	// Update the tender in the database
	query := `UPDATE tender SET name=$1, description=$2, service_type=$3, status=$4, version=$5, updated_at=$6, creator_username=$7
              WHERE id=$8 RETURNING id`
	var tenderID uuid.UUID
	if err := db.QueryRow(query, updates.Name, updates.Description, updates.ServiceType, updates.Status, current.Version, updates.UpdatedAt, user, id).Scan(&tenderID); err != nil {
		return uuid.Nil, fmt.Errorf("failed to update tender: %v", err)
	}

	return tenderID, nil
}
