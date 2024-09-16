package database

import (
	"database/sql"
	"github.com/google/uuid"
	"src/models"
)

func BidRollBack(db *sql.DB, version int, id uuid.UUID, bid *models.Bid) error {
	current := models.Bid{}

	err := getBidData(db, id, &current)
	if err != nil {
		return err
	}

	err = saveCurrentBidToHistory(db, &current)
	if err != nil {
		return err
	}

	err = GetBidDataFromHistory(db, id, bid, version)
	if err != nil {
		return err
	}

	bid.Version = current.Version + 1

	return nil

}
