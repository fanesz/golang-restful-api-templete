package database

import (
	"backend/app/api/user"

	"gorm.io/gorm"
)

type MigratableModel interface {
	Migrate(db *gorm.DB) error
}

var modelList = []interface{}{
	&user.User{},
	&user.Role{},
}

func Migrate(db *gorm.DB) error {
	for _, model := range modelList {
		if err := db.AutoMigrate(model); err != nil {
			panic(err)
		}
	}
	return nil
}
