package repository

import (
	"SWIFT_task/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.AutoMigrate(&model.Country{})
	if err != nil {
		log.Fatalf("Failed to create tabels in database: %v", err)
	}
	err = db.AutoMigrate(&model.Bank{})
	if err != nil {
		log.Fatalf("Failed to create tabels in database: %v", err)
	}
	err = db.AutoMigrate(&model.BankRelationship{})
	if err != nil {
		log.Fatalf("Failed to create tabels in database: %v", err)
	}

	branch1 := &model.Bank{
		SwiftCode:   "AIPOPLP1XXX",
		Address:     "123 Test St",
		BankName:    "Test Bank1",
		CountryISO2: "PL",
	}

	branch2 := &model.Bank{
		SwiftCode:   "AIPOPLP1FGD",
		Address:     "456 Another St",
		BankName:    "Another Bank",
		CountryISO2: "PL",
	}

	branch3 := &model.Bank{
		SwiftCode:   "FEDCBA12XXX",
		Address:     "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",
		BankName:    "UNITED BANK OF ALBANIA SH.A",
		CountryISO2: "BG",
	}

	country1 := &model.Country{
		CountryISO2: "PL",
		CountryName: "POLAND",
	}

	country2 := &model.Country{
		CountryISO2: "BG",
		CountryName: "BULGARIA",
	}

	db.Create(branch1)
	db.Create(branch2)
	db.Create(branch3)

	db.Create(&model.BankRelationship{
		HeadquarterSwiftCode: "AIPOPLP1XXX",
		BranchSwiftCode:      "AIPOPLP1FGD",
		Headquarter:          branch1,
		Branch:               branch2,
	})

	db.Create(&country1)
	db.Create(&country2)

	return db
}

func TestBankRepository_GetBankBySwiftCode(t *testing.T) {
	db := setupTestDB()
	bankRepository := NewBankRepository(db)

	t.Run("should get bank successfully", func(t *testing.T) {
		// given
		bank1 := &model.Bank{
			SwiftCode:   "AIPOPLP1XXX",
			Address:     "123 Test St",
			BankName:    "Test Bank1",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		swiftCode1 := bank1.SwiftCode

		bank2 := &model.Bank{
			SwiftCode:   "AIPOPLP1FGD",
			Address:     "456 Another St",
			BankName:    "Another Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		swiftCode2 := bank2.SwiftCode

		// when
		resultBank1, err1 := bankRepository.GetBankBySwiftCode(swiftCode1)
		resultBank2, err2 := bankRepository.GetBankBySwiftCode(swiftCode2)

		// then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, bank1, resultBank1)
		assert.Equal(t, bank2, resultBank2)
	})

	t.Run("should return nil when bank with given swift-code does not exist", func(t *testing.T) {
		// given
		swiftCode1 := "BIPOPLP1XXX"
		swiftCode2 := "BIPOPLP1ABC"

		// when
		resultBank1, err1 := bankRepository.GetBankBySwiftCode(swiftCode1)
		resultBank2, err2 := bankRepository.GetBankBySwiftCode(swiftCode2)

		// then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Nil(t, resultBank1)
		assert.Nil(t, resultBank2)

	})

}

func TestBankRepository_GetBranchesForHeadquarter(t *testing.T) {
	db := setupTestDB()
	bankRepository := NewBankRepository(db)

	t.Run("should return branches properly for headquarter", func(t *testing.T) {
		// given
		headquarterSwift := "AIPOPLP1XXX"
		expectedBranches := []model.Bank{
			{
				SwiftCode:   "AIPOPLP1FGD",
				Address:     "456 Another St",
				BankName:    "Another Bank",
				CountryISO2: "PL",
				Country: model.Country{
					CountryISO2: "PL",
					CountryName: "POLAND",
				},
			},
		}

		// when
		resultBranches, err := bankRepository.GetBranchesForHeadquarter(headquarterSwift)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedBranches, resultBranches)
	})

	t.Run("should return empty list because headquarter does not have branches", func(t *testing.T) {
		// given
		headquarterSwift := "FEDCBA12XXX"

		// when
		resultBranches, err := bankRepository.GetBranchesForHeadquarter(headquarterSwift)

		// then
		assert.NoError(t, err)
		assert.Equal(t, []model.Bank{}, resultBranches)
	})

}

func TestBankRepository_GetBanksByISO2code(t *testing.T) {
	db := setupTestDB()
	bankRepository := NewBankRepository(db)

	t.Run("should get two banks successfully", func(t *testing.T) {
		// given
		isoCodeWithTwoBanks := "PL"

		twoBanks := []model.Bank{
			{
				SwiftCode:   "AIPOPLP1XXX",
				Address:     "123 Test St",
				BankName:    "Test Bank1",
				CountryISO2: "PL",
				Country: model.Country{
					CountryISO2: "PL",
					CountryName: "POLAND",
				},
			},
			{
				SwiftCode:   "AIPOPLP1FGD",
				Address:     "456 Another St",
				BankName:    "Another Bank",
				CountryISO2: "PL",
				Country: model.Country{
					CountryISO2: "PL",
					CountryName: "POLAND",
				},
			},
		}

		// when
		banks, err := bankRepository.GetBanksByISO2code(isoCodeWithTwoBanks)

		// then
		assert.NoError(t, err)
		assert.Equal(t, twoBanks, banks)

	})

	t.Run("should return one bank in a list successfully", func(t *testing.T) {
		// given
		isoCodeWithOneBank := "BG"
		oneBank := []model.Bank{
			{
				SwiftCode:   "FEDCBA12XXX",
				Address:     "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",
				BankName:    "UNITED BANK OF ALBANIA SH.A",
				CountryISO2: "BG",
				Country: model.Country{
					CountryISO2: "BG",
					CountryName: "BULGARIA",
				},
			},
		}

		// when
		banks, err := bankRepository.GetBanksByISO2code(isoCodeWithOneBank)

		// then
		assert.NoError(t, err)
		assert.Equal(t, oneBank, banks)
	})

	t.Run("should return empty list because there is no banks with such ISO2 code", func(t *testing.T) {
		// given
		isoCodeWithoutBanks := "AB"

		// when
		banks, err := bankRepository.GetBanksByISO2code(isoCodeWithoutBanks)

		// then
		assert.NoError(t, err)
		assert.Equal(t, []model.Bank{}, banks)
	})
}

func TestBankRepository_AddBank(t *testing.T) {
	db := setupTestDB()
	bankRepository := NewBankRepository(db)

	t.Run("should create bank successfully", func(t *testing.T) {
		// given
		bank := &model.Bank{
			SwiftCode:   "ABCDEF12XXX",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
		}

		// when
		err := bankRepository.AddBank(bank)

		// then
		assert.NoError(t, err)

		var storedBank model.Bank
		err = db.First(&storedBank, "swift_code = ?", bank.SwiftCode).Error
		assert.NoError(t, err)
		assert.Equal(t, bank, &storedBank)
	})

	t.Run("should return error when creating bank with duplicate SwiftCode", func(t *testing.T) {
		// given
		duplicatedBank := &model.Bank{
			SwiftCode:   "ABCDEF12XXX",
			Address:     "Another Address",
			BankName:    "Another Bank",
			CountryISO2: "DE",
		}

		// when
		err := bankRepository.AddBank(duplicatedBank)

		// then
		assert.Error(t, err)
	})
}

func TestBankRepository_RemoveBranchBySwiftCode(t *testing.T) {
	db := setupTestDB()
	bankRepository := NewBankRepository(db)

	t.Run("should remove bank successfully", func(t *testing.T) {
		// given
		bank1 := &model.Bank{
			SwiftCode:   "AIPOPLP1XXX",
			Address:     "123 Test St",
			BankName:    "Test Bank1",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}

		bank2 := &model.Bank{
			SwiftCode:   "AIPOPLP1FGD",
			Address:     "456 Another St",
			BankName:    "Another Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}

		// when & then
		resultBank, _ := bankRepository.GetBankBySwiftCode(bank1.SwiftCode)
		assert.Equal(t, bank1, resultBank)
		err := bankRepository.RemoveBankBySwiftCode(bank1.SwiftCode)
		assert.NoError(t, err)
		resultBank, _ = bankRepository.GetBankBySwiftCode(bank1.SwiftCode)
		assert.Nil(t, resultBank)

		resultBank2, _ := bankRepository.GetBankBySwiftCode(bank2.SwiftCode)
		assert.Equal(t, bank2, resultBank2)
		err = bankRepository.RemoveBankBySwiftCode(bank2.SwiftCode)
		assert.NoError(t, err)
		resultBank2, _ = bankRepository.GetBankBySwiftCode(bank2.SwiftCode)
		assert.Nil(t, resultBank2)
	})

}
