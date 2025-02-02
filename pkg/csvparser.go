package pkg

import (
	"SWIFT_task/internal/model"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func ParseCsvRows(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprintf("Error closing file: %v", err))
		}
	}(file)

	reader := csv.NewReader(file)
	return readCsvWithoutHeaders(reader)
}

func GetBanksFromRecords(records [][]string) map[string]model.Bank {
	banks := make(map[string]model.Bank, len(records))

	for _, record := range records {
		address := record[4]
		bankName := record[3]
		countryISO2 := strings.ToUpper(record[0])
		countryName := strings.ToUpper(record[6])
		swiftCode := record[1]

		bankRecord := model.Bank{
			Address:     address,
			BankName:    bankName,
			CountryISO2: countryISO2,
			SwiftCode:   swiftCode,
			Country: model.Country{
				CountryISO2: countryISO2,
				CountryName: countryName,
			},
		}

		banks[swiftCode] = bankRecord
	}

	return banks
}

func GetRelationshipsFromBanks(banks map[string]model.Bank) []model.BankRelationship {
	relationships := make([]model.BankRelationship, 0)

	for swiftCode, bank := range banks {
		if !bank.IsHeadquarter() {
			headquarterSwiftCode := getHeadquarterSwiftCode(swiftCode)
			headquarter, ok := banks[headquarterSwiftCode]
			if ok {
				relationship := model.BankRelationship{
					HeadquarterSwiftCode: headquarterSwiftCode,
					BranchSwiftCode:      swiftCode,
					Headquarter:          &headquarter,
					Branch:               &bank,
				}
				relationships = append(relationships, relationship)
			}
		}
	}
	return relationships
}

func GetCountriesFromRecords(records [][]string) []model.Country {
	countryMap := make(map[string]model.Country)

	for _, record := range records {
		countryISO2 := strings.ToUpper(record[0])
		countryName := strings.ToUpper(record[6])

		if _, exists := countryMap[countryISO2]; !exists {
			countryMap[countryISO2] = model.Country{
				CountryISO2: countryISO2,
				CountryName: countryName,
			}
		}
	}

	countries := make([]model.Country, len(countryMap))
	index := 0
	for _, country := range countryMap {
		countries[index] = country
		index++
	}

	return countries
}

func getHeadquarterSwiftCode(swiftCode string) string {
	return swiftCode[:8] + "XXX"
}

func readCsvWithoutHeaders(reader *csv.Reader) ([][]string, error) {
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
