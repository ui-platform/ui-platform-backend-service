package entity

import "time"

const (
	ScreenStatusPublished   = "Published"
	ScreenStatusUnPublished = "Unpublished"
	ScreenStatusArchived    = "Archived"
)

type Screen struct {
	Id          string                 `json:"id" db:"id"`
	ProjectId   string                 `json:"project_id" db:"project_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description,omitempty" db:"description"`
	Status      string                 `json:"status,omitempty" db:"status"`
	Widgets     map[string]interface{} `json:"widgets,omitempty" db:"widgets"`
	Settings    map[string]interface{} `json:"settings,omitempty" db:"settings"`
	CreatedAt   time.Time              `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   time.Time              `json:"deleted_at,omitempty" db:"deleted_at"`
}
