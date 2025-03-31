package entity

import (
	"errors"
	"time"
)

const (
	ProjectStatusPublished   = "Published"
	ProjectStatusUnPublished = "Unpublished"
	ProjectStatusArchived    = "Archived"
)

type Project struct {
	ID          string    `json:"id,omitempty" db:"id"`
	Name        string    `json:"name,omitempty" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Status      string    `json:"status,omitempty" db:"status"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (p *Project) EntityName() string {
	return "projects"
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if len(p.Name) > 100 {
		return errors.New("name must be less than 100 characters")
	}

	if p.Description == "" {
		return errors.New("description is required")
	}
	if len(p.Description) > 4000 {
		return errors.New("description must be less than 1000 characters")
	}

	return nil
}
