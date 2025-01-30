package model

type Branch struct {
	Address       string `gorm:"size:255"`
	BankName      string `gorm:"size:255"`
	CountryISO2   string `gorm:"size:2"`
	CountryName   string `gorm:"size:255"`
	SwiftCode     string `gorm:"size:11;primaryKey"`
	IsHeadquarter bool   `gorm:"default:false"`
}

type BranchRelationship struct {
	HeadquarterSwiftCode    string `gorm:"size:11;primaryKey"`
	OrdinaryBranchSwiftCode string `gorm:"size:11;primaryKey"`

	HeadquarterBranch *Branch `gorm:"foreignKey:HeadquarterSwiftCode;references:SwiftCode"`
	OrdinaryBranch    *Branch `gorm:"foreignKey:OrdinaryBranchSwiftCode;references:SwiftCode"`
}

func (b Branch) GetHeadQuarterSwiftCode() string {
	return b.SwiftCode[:8] + "XXX"
}
