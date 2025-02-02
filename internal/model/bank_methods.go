package model

import "strings"

func (b Bank) GetHeadQuarterSwiftCode() string {
	return b.SwiftCode[:8] + "XXX"
}

func (b Bank) ToBankWithoutCountryNameDto() BankWithoutCountryNameDto {
	return BankWithoutCountryNameDto{
		Address:       b.Address,
		BankName:      b.BankName,
		CountryISO2:   b.CountryISO2,
		IsHeadquarter: b.IsHeadquarter(),
		SwiftCode:     b.SwiftCode,
	}
}

func (b Bank) ToBankDto() BankDto {
	return BankDto{
		Address:       b.Address,
		BankName:      b.BankName,
		CountryISO2:   b.CountryISO2,
		CountryName:   b.Country.CountryName,
		IsHeadquarter: b.IsHeadquarter(),
		SwiftCode:     b.SwiftCode,
		Branches:      nil,
	}

}

func (b Bank) IsHeadquarter() bool {
	return strings.HasSuffix(b.SwiftCode, "XXX")
}
