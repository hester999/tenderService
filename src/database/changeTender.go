package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"src/models"
)

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
	// SQL запрос для получения данных о тендере
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
		&tender.,
		&tender.CreatorUsername,
		&tender.Version,
		&tender.CreatedAt,
	)

	// Обработка возможной ошибки при выполнении запроса
	if err != nil {
		return fmt.Errorf("error fetching tender data: %v", err)
	}

	return nil
}


func ChangeTender(db *sql.DB, id uuid.UUID, tender *models.Tender, resp *models.TenderResponse) error {
	tempTender := models.Tender{}
	err := checkTender(db, id)
	if err != false {
		return fmt.Errorf("tender not  found")
	}
	getTenderData(db, id,&tempTender)
	return nil
}
