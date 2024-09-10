package database

import (
	"database/sql"
	"src/models"
)

func ChangeTenderStatus(db *sql.DB, response *models.TenderResponse, status, user string) error {

	//err := ValidatePermission()
	return nil
}
