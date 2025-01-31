package model

type Branch struct {
	Address       string `gorm:"size:255" json:"address"`
	BankName      string `gorm:"size:255" json:"bankName"`
	CountryISO2   string `gorm:"size:2" json:"countryISO2"`
	CountryName   string `gorm:"size:255" json:"countryName"`
	IsHeadquarter bool   `gorm:"default:false" json:"isHeadquarter"`
	SwiftCode     string `gorm:"size:11;primaryKey" json:"swiftCode"`
}

type BranchRelationship struct {
	HeadquarterSwiftCode    string `gorm:"size:11;primaryKey"`
	OrdinaryBranchSwiftCode string `gorm:"size:11;primaryKey"`

	HeadquarterBranch *Branch `gorm:"foreignKey:HeadquarterSwiftCode;references:SwiftCode"`
	OrdinaryBranch    *Branch `gorm:"foreignKey:OrdinaryBranchSwiftCode;references:SwiftCode"`
}

type BranchWithoutCountryNameDto struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type BranchDto struct {
	Address       string                        `json:"address"`
	BankName      string                        `json:"bankName"`
	CountryISO2   string                        `json:"countryISO2"`
	CountryName   string                        `json:"countryName"`
	IsHeadquarter bool                          `json:"isHeadquarter"`
	SwiftCode     string                        `json:"swiftCode"`
	Branches      []BranchWithoutCountryNameDto `json:"branches,omitempty"`
}

type BranchesForCountryDto struct {
	CountryISO2 string                        `json:"countryISO2"`
	CountryName string                        `json:"countryName"`
	Branches    []BranchWithoutCountryNameDto `json:"swiftCodes"`
}

func (b Branch) GetHeadQuarterSwiftCode() string {
	return b.SwiftCode[:8] + "XXX"
}

func (b Branch) ToBranchWithoutCountryNameDto() BranchWithoutCountryNameDto {
	return BranchWithoutCountryNameDto{
		Address:       b.Address,
		BankName:      b.BankName,
		CountryISO2:   b.CountryISO2,
		IsHeadquarter: b.IsHeadquarter,
		SwiftCode:     b.SwiftCode,
	}
}

func (b Branch) ToBranchDto() BranchDto {
	return BranchDto{
		Address:       b.Address,
		BankName:      b.BankName,
		CountryISO2:   b.CountryISO2,
		CountryName:   b.CountryName,
		IsHeadquarter: b.IsHeadquarter,
		SwiftCode:     b.SwiftCode,
		Branches:      nil,
	}

}
