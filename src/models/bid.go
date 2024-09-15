package models

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"src/internal"
	"time"
)

type Bid struct {
	Id             uuid.UUID `json:"id"`
	TenderId       uuid.UUID `json:"tenderId"`
	OrganizationId uuid.UUID `json:"organizationId"`
	CreatorUserId  uuid.UUID `json:"creatorUserId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"` // Включать в JSON пустое поле
	Status         string    `json:"status"`
	AuthorType     string    `json:"authorType"`
	AuthorId       uuid.UUID `json:"authorId"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type BidResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	AuthorType string `json:"authorType"`
	AuthorId   string `json:"authorId"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

func (b *Bid) ValidateBidUser(db *sql.DB) error {
	err := internal.ValidateUserBid(db, b.CreatorUserId)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bid) ValidateTenderId(db *sql.DB) error {
	err := internal.ValidateTenderId(db, b.TenderId)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bid) ValidateFields() error {
	if b.Name == "" {
		return errors.New("name cannot be empty")
	}
	if b.Description == "" {
		return errors.New("description cannot be empty")
	}
	if b.TenderId == uuid.Nil {
		return errors.New("tenderId cannot be empty")
	}
	if b.AuthorId == uuid.Nil {
		return errors.New("authorId cannot be empty")
	}

	if b.Status == "" {
		b.Status = "Created"
	}
	if b.CreatorUserId == uuid.Nil {
		b.CreatorUserId = b.AuthorId
	}
	b.Version = 1
	return nil
}

type ResponseGetBid struct {
	Bid    []Bid `json:"bids"`
	Offset int
	Limit  int
	User   string
}

func (r *ResponseGetBid) ValidateUser(user string, db *sql.DB) error {
	var exist bool
	err := db.QueryRow(`select exists(select 1 from employee where  username = $1)`, user).Scan(&exist)

	if err != nil {
		return err
	}

	if !exist {
		return errors.New("user not fond ")
	}

	return nil
}

func (r *ResponseGetBid) ValidatePermission(user string, db *sql.DB) error {
	var exist bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM organization_responsible WHERE user_id = (SELECT id FROM employee WHERE username = $1))`, user).Scan(&exist)

	if err != nil {
		return err
	}
	if !exist {
		return errors.New("user dont have permission")

	}
	return err
}
