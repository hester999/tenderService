package internal

import (
	"database/sql"
	"github.com/google/uuid"
)

func GetVersion(db *sql.DB, id uuid.UUID) int {

	var version int

	_ = db.QueryRow(`select version from tender where id = $1`, id).Scan(&version)

	return version
}

func GetVersionFromHistoryTender(db *sql.DB, id uuid.UUID, vers int) bool {
	query := "select exists(select tender_id,version from tender_history where tender_id = $1 and version = $2)"

	var status bool
	_ = db.QueryRow(query, id, vers).Scan(&status)
	if status == false {

		return status
	}
	return true
}
