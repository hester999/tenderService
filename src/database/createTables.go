package database

import (
	"database/sql"
	"fmt"
)

func createTables(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	tableCreationQueries := []string{

		`CREATE TABLE IF NOT EXISTS organization (
		id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name        VARCHAR(100) NOT NULL,
		description TEXT,
		type        organization_type,
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS employee (
		id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		username   VARCHAR(50) UNIQUE NOT NULL,
		first_name VARCHAR(50),
		last_name  VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS organization_responsible (
		id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		organization_id UUID REFERENCES organization ON DELETE CASCADE,
		user_id         UUID REFERENCES employee ON DELETE CASCADE
	);`,

		`CREATE TABLE IF NOT EXISTS tender (
		id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name             VARCHAR(255) NOT NULL,
		description      TEXT,
		service_type     VARCHAR(100),
		status           VARCHAR(50) DEFAULT 'Created', 
		organization_id  UUID REFERENCES organization ON DELETE CASCADE,
		creator_user_id  UUID REFERENCES employee ON DELETE CASCADE,
		creator_username VARCHAR(50) REFERENCES employee(username),
		version          INTEGER DEFAULT 1,
		created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS tender_history (
		history_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		tender_id        UUID REFERENCES tender ON DELETE CASCADE,
		name             VARCHAR(255),
		description      TEXT,
		service_type     VARCHAR(100),
		status           VARCHAR(50),
		version          INTEGER NOT NULL,
		organization_id  UUID,
		creator_user_id  UUID,
		creator_username VARCHAR(50),
		created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tender_id, version)
	);`,

		`CREATE TABLE IF NOT EXISTS bid (
		id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		tender_id       UUID REFERENCES tender ON DELETE CASCADE,
		creator_user_id UUID REFERENCES employee ON DELETE CASCADE,
		name            VARCHAR(255) NOT NULL,
		description     TEXT,
		status          VARCHAR(50) DEFAULT 'Created',  
		author_type     VARCHAR(50), 
		author_id       UUID,        
		version         INTEGER DEFAULT 1,
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tender_id, version)
	);`,

		`CREATE TABLE IF NOT EXISTS bid_history (
		history_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		bid_id          UUID REFERENCES bid ON DELETE CASCADE,
		tender_id       UUID REFERENCES tender ON DELETE SET NULL,
		creator_user_id UUID REFERENCES employee ON DELETE SET NULL,
		name            VARCHAR(255),
		description     TEXT,
		service_type    VARCHAR(100),
		status          VARCHAR(50), 
		author_type     VARCHAR(50), 
		author_id       UUID,        
		version         INTEGER NOT NULL,
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (bid_id, version)
	);`,
	}

	for _, query := range tableCreationQueries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute query: %v, error: %w", query, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
