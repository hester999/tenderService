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
