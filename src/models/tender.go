package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Tender struct {
	Id              uuid.UUID `json:"id"` // Изменено на UUID
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	ServiceType     string    `json:"serviceType"`
	Version         int       `json:"version"`
	CreatedAt       string    `json:"createdAt"`
	OrganizationId  uuid.UUID `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
}

func (t *Tender) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if t.Description == "" {
		return fmt.Errorf("description is required")
	}
	if t.Status == "" {
		return fmt.Errorf("status is required")
	}
	if t.ServiceType == "" {
		return fmt.Errorf("serviceType is required")
	}
	if t.OrganizationId == uuid.Nil { // Исправлено: проверка на пустой UUID
		return fmt.Errorf("organizationId is required")
	}
	if t.CreatorUsername == "" {
		return fmt.Errorf("creatorUsername is required")
	}
	return nil
}

type ResponseGetBanner struct {
	Tenders []Tender `json:"tenders"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	User    string   `json:"user"`
}

func (t *ResponseGetBanner) Validatepaging() {
	if t.Limit > 5 {
		t.Limit = 5
	}
}

type TenderResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ServiceType string `json:"serviceType"`
	Version     int    `json:"version"` // Исправлено: правильно написано "version"
	CreatedAt   string `json:"createdAt"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}
