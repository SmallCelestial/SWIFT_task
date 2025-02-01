package db

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/pkg"
	"gorm.io/gorm"
	"log"
)

func SaveData(fileName string, db *gorm.DB) {
	records, err := pkg.ParseCsvRows(fileName)
	if err != nil {
		log.Fatal(err)
	}

	countries := pkg.GetCountriesFromRecords(records)
	banks := pkg.GetBanksFromRecords(records)
	relationships := pkg.GetRelationshipsFromBanks(banks)

	saveCountries(countries, db)
	saveBanks(banks, db)
	saveRelationships(relationships, db)

}

func saveBanks(banks map[string]model.Bank, db *gorm.DB) {
	for _, bank := range banks {
		result := db.Create(&bank)
		if result.Error != nil {
			log.Printf("Error saving banks %s", result.Error.Error())
		}
	}
}

func saveRelationships(relationships []model.BankRelationship, db *gorm.DB) {
	for _, relationship := range relationships {
		result := db.Create(&relationship)
		if result.Error != nil {
			log.Printf("Error saving relationships %s", result.Error.Error())
		}
	}
}

func saveCountries(countries []model.Country, db *gorm.DB) {
	for _, country := range countries {
		result := db.Create(&country)
		if result.Error != nil {
			log.Printf("Error saving countries %s", result.Error.Error())
		}
	}
}
