package database

import (
	"backend/app/common/utils"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) {
	InitRole(db)
	InitAdmin(db)
}

func InitRole(db *gorm.DB) {
	var isExists bool
	query := db.Table("roles").Select("count(*) > 0").Find(&isExists)
	if query.Error != nil {
		return
	}

	if !isExists {
		adminRoleID := utils.GenerateRandomStringWithNumber(24)
		userRoleID := utils.GenerateRandomStringWithNumber(24)
		guestRoleID := utils.GenerateRandomStringWithNumber(24)

		db.Exec(`INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, 'ADMIN', now(), now())`, adminRoleID)
		db.Exec(`INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, 'USER', now(), now())`, userRoleID)
		db.Exec(`INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, 'GUEST', now(), now())`, guestRoleID)
	}
}

func InitAdmin(db *gorm.DB) {
	email := utils.GetEnv("USER_ADMIN_EMAIL")

	var isExists bool
	query := db.Table("users").Select("count(*) > 0").Where("email = ?", email).Find(&isExists)
	if query.Error != nil {
		return
	}

	if !isExists {
		uuid := utils.GenerateUUID()
		password := utils.GetEnv("USER_ADMIN_PASSWORD")
		utils.Encrypt(&password)

		var roleID string
		db.Table("roles").Select("id").Where("name = 'ADMIN'").Find(&roleID)
		db.Exec(`INSERT INTO users (uuid, username, email, password, role_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`, uuid, email, email, password, roleID)
	}
}
