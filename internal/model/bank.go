package model

type Bank struct {
	Address     string `gorm:"size:255"`
	BankName    string `gorm:"size:255"`
	CountryISO2 string `gorm:"size:2"`
	SwiftCode   string `gorm:"size:11;primaryKey"`

	Country Country `gorm:"foreignKey:CountryISO2;references:CountryISO2"`
}

type BankRelationship struct {
	HeadquarterSwiftCode string `gorm:"size:11;primaryKey"`
	BranchSwiftCode      string `gorm:"size:11;primaryKey"`

	Headquarter *Bank `gorm:"foreignKey:HeadquarterSwiftCode;references:SwiftCode"`
	Branch      *Bank `gorm:"foreignKey:BranchSwiftCode;references:SwiftCode"`
}

type Country struct {
	CountryISO2 string `gorm:"size:2;primaryKey"`
	CountryName string `gorm:"size:255"`
}
