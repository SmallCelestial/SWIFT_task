package pkg

import (
	"SWIFT_task/internal/model"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func ParseSwiftCsvToBranches(filePath string) (map[string]*model.Branch, error) {
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
	records, err := readCsvWithoutHeaders(reader)

	if err != nil {
		return nil, err
	}

	branches := getBranchesFromRecords(records)
	return branches, nil
}

func isHeadQuarter(swiftCode string) bool {
	return strings.HasSuffix(swiftCode, "XXX")
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

func getBranchesFromRecords(records [][]string) map[string]*model.Branch {
	branches := make(map[string]*model.Branch, len(records))

	for _, record := range records {
		address := record[4]
		bankName := record[3]
		countryISO2 := strings.ToUpper(record[0])
		countryName := strings.ToUpper(record[6])
		swiftCode := record[1]
		isHeadquarter := isHeadQuarter(swiftCode)

		branchRecord := &model.Branch{
			Address:       address,
			BankName:      bankName,
			CountryISO2:   countryISO2,
			CountryName:   countryName,
			SwiftCode:     swiftCode,
			IsHeadquarter: isHeadquarter,
		}

		branches[swiftCode] = branchRecord
	}

	return branches
}
