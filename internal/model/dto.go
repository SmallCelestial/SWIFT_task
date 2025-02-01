package model

type BankWithoutCountryNameDto struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"IsHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type BankDto struct {
	Address       string                      `json:"address"`
	BankName      string                      `json:"bankName"`
	CountryISO2   string                      `json:"countryISO2"`
	CountryName   string                      `json:"countryName"`
	IsHeadquarter bool                        `json:"IsHeadquarter"`
	SwiftCode     string                      `json:"swiftCode"`
	Branches      []BankWithoutCountryNameDto `json:"branches,omitempty"`
}

type CountryBanksDto struct {
	CountryISO2 string                      `json:"countryISO2"`
	CountryName string                      `json:"countryName"`
	Branches    []BankWithoutCountryNameDto `json:"swiftCodes"`
}
