package pkg

import (
	"SWIFT_task/internal/model"
	"encoding/csv"
	"os"
	"strings"
)

func ParseSwiftCsvToBranches(filePath string) (map[string]*model.Branch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := readCsvWithoutHeaders(reader)

	if err != nil {
		return nil, err
	}

	branches := getBranchesFromRecords(records)
	assignBranchesToHeadquarters(branches)
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

func assignBranchesToHeadquarters(branches map[string]*model.Branch) {
	for _, branch := range branches {
		if !branch.IsHeadquarter {
			headquarter := findHeadquarterForBranch(branch, branches)
			if headquarter != nil {
				headquarter.Branches = append(headquarter.Branches, *branch)
			}
		}
	}
}

func findHeadquarterForBranch(branch *model.Branch, branches map[string]*model.Branch) *model.Branch {
	headquarterSwiftCode := branch.SwiftCode[:8] + "XXX"
	headquarter, ok := branches[headquarterSwiftCode]
	if !ok {
		return nil
	}
	return headquarter
}
