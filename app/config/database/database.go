package database

import (
	"backend/app/common/utils"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseInstance *gorm.DB

func InitializeDB() {
	fmt.Println("===== Initialize Database =====")
	db, err := connectDB()
	Seeder(db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("connected to: %s\n", utils.GetEnv("DB_NAME"))

	databaseInstance = db
}

func GetDBInstance() *gorm.DB {
	return databaseInstance
}

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		utils.GetEnv("DB_HOST"),
		utils.GetEnv("DB_PORT"),
		utils.GetEnv("DB_USER"),
		utils.GetEnv("DB_PASSWORD"),
		utils.GetEnv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := Migrate(db); err != nil {
		fmt.Printf("error migrate: %v", err)
	}
	return db, nil
}
