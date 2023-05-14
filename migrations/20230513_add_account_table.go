package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username       string    `gorm:"uniqueIndex"`
	HashedPassword string
	gorm.Model
}

func CreateAccountTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202305130001",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(Account{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(Account{})
		},
	}
}
