package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity interface {
	GetId() string
}

type EntityBase struct {
	ID        string          `gorm:"type:uuid;primaryKey;column:id"`
	CreatedAt time.Time       `gorm:"column:created_at;<-:update"`
	UpdatedAt time.Time       `gorm:"column:updated_at;<-:update"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (e *EntityBase) GetId() string {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}

	return e.ID
}

func (u *EntityBase) BeforeCreate(tx *gorm.DB) error {
	u.GetId()

	u.DeletedAt = nil
	return nil
}
