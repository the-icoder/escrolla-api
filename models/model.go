package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        string         `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt int64          `json:"created_at,omitempty"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
