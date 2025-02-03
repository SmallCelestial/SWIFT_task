package db

import (
	"SWIFT_task/config"
	"SWIFT_task/internal/model"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDB() *gorm.DB {
	db := config.InitDatabaseConnection()

	err := db.AutoMigrate(&model.Country{})
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
		filePath := os.Getenv("DATA_FILE_PATH")
		if filePath == "" {
			filePath = "data/swift-codes.csv"
		}
		SaveData("/app/data/swift-codes.csv", db)
	}
}
