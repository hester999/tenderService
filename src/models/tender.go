package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Tender struct {
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	ServiceType     string    `json:"serviceType"`
	Version         int       `json:"version"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string
	OrganizationId  uuid.UUID `json:"organizationId"`
	CreatorUserId   uuid.UUID
	CreatorUsername string `json:"creatorUsername"`
}

func (t *Tender) ValidateChangeFiled() error {
	if t.Version != 0 || t.CreatedAt != "" || t.OrganizationId != uuid.Nil || t.CreatorUsername != "" {
		fields := []string{}
		if t.Version != 0 {
			fields = append(fields, "Version")
		}
		if t.CreatedAt != "" {
			fields = append(fields, "CreatedAt")
		}
		if t.OrganizationId != uuid.Nil {
			fields = append(fields, "OrganizationId")
		}
		if t.CreatorUsername != "" {
			fields = append(fields, "CreatorUsername")
		}
		return fmt.Errorf("these fields cannot be modified: %v", fields)
	}
	return nil

}

func (t *Tender) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if t.Description == "" {
		return fmt.Errorf("description is required")
	}
	if t.Status == "" {
		t.Status = "Created"
	}
	if t.ServiceType == "" {
		return fmt.Errorf("serviceType is required")
	}
	if t.OrganizationId == uuid.Nil {
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
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}
