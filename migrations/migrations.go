package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RunMigration(gormDB *gorm.DB) {
	m := gormigrate.New(gormDB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		CreateUserTable(),
	})
	if err := m.Migrate(); err != nil {
		log.Fatal().Err(err).Msgf("migration failed: %v", err.Error())
	}
}
