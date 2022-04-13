package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/**
 * @model Base
 * @description This is injected into other models to provide common functionality.
 */

type Base struct {
	ID        string `gorm:"primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.NewString()
	return
}
