package migrate

import (
	"github.com/fuvy/effmob-test/internal/storage"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}
	return db.AutoMigrate(
		&storage.Car{},
	)
}
