package database

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"src/models"
)

func TenderRollBack(db *sql.DB, id uuid.UUID, user string, tender *models.Tender, version int) error {
	current := models.Tender{}
	if !checkTender(db, id) {
		return errors.New("TenderRollback fail")
	}

	err := getTenderData(db, id, &current)
	if err != nil {
		return err
	}

	err = saveCurrentTenderToHistory(db, &current)
	if err != nil {
		return err
	}

	err = GetDataFromHistoryTender(db, id, tender, version)
	if err != nil {
		return err
	}

	tender.Version = current.Version + 1

	return nil
}
