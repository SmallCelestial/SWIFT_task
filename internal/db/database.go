package db

import (
	"SWIFT_task/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	//host := os.Getenv("DB_HOST")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_PASSWORD")
	//dbname := os.Getenv("DB_NAME")
	//port := "5432"
	host := "localhost"
	user := "postgres"
	password := "admin"
	dbname := "swift_task"
	port := "5432"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&model.Country{})
	if err != nil {
		log.Fatalf("Failed to migrate Country table: %v", err)
	}

	err = db.AutoMigrate(&model.Bank{})
	if err != nil {
		log.Fatalf("Failed to migrate Bank table: %v", err)
	}

	err = db.AutoMigrate(&model.BankRelationship{})
	if err != nil {
		log.Fatalf("Failed to migrate BankRelationship table: %v", err)
	}

	fillData(db)

	return db
}

func fillData(db *gorm.DB) {
	var bankCount int64
	var relationshipCount int64
	var countryCount int64
	db.Model(&model.Bank{}).Count(&bankCount)
	db.Model(&model.BankRelationship{}).Count(&relationshipCount)
	db.Model(&model.Country{}).Count(&countryCount)
	if bankCount == 0 && relationshipCount == 0 && countryCount == 0 {
		SaveData("data/swift-codes.csv", db)
	}

}
