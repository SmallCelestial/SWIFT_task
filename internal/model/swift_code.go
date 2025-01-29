package model

type Branch struct {
	Address       string   `json:"address"`
	BankName      string   `json:"bankName"`
	CountryISO2   string   `json:"countryISO2"`
	CountryName   string   `json:"countryName"`
	SwiftCode     string   `json:"swiftCode"`
	IsHeadquarter bool     `json:"isHeadquarter"`
	Branches      []Branch `json:"branches,omitempty"`
}
