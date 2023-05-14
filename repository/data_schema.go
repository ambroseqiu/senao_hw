package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username       string    `gorm:"uniqueIndex"`
	HashedPassword string
	gorm.Model
}
