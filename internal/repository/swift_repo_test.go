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

	err = db.AutoMigrate(&model.Branch{})
	err = db.AutoMigrate(&model.BranchRelationship{})
	if err != nil {
		log.Fatalf("Failed to create tabels in database: %v", err)
	}

	branch1 := &model.Branch{
		SwiftCode:     "AIPOPLP1XXX",
		Address:       "123 Test St",
		BankName:      "Test Bank1",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}

	branch2 := &model.Branch{
		SwiftCode:     "AIPOPLP1FGD",
		Address:       "456 Another St",
		BankName:      "Another Bank",
		CountryISO2:   "PL",
		CountryName:   "Germany",
		IsHeadquarter: false,
	}

	branch3 := &model.Branch{
		SwiftCode:     "FEDCBA12XXX",
		Address:       "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",
		BankName:      "UNITED BANK OF ALBANIA SH.A",
		CountryISO2:   "BG",
		CountryName:   "Germany",
		IsHeadquarter: true,
	}

	db.Create(branch1)
	db.Create(branch2)
	db.Create(branch3)

	db.Create(&model.BranchRelationship{
		HeadquarterSwiftCode:    "AIPOPLP1XXX",
		OrdinaryBranchSwiftCode: "AIPOPLP1FGD",
		HeadquarterBranch:       branch1,
		OrdinaryBranch:          branch2,
	})

	return db
}

func TestBranchRepository_CreateBranch(t *testing.T) {
	db := setupTestDB()
	branchRepository := NewBranchRepository(db)

	t.Run("should create branch successfully", func(t *testing.T) {
		// given
		branch := &model.Branch{
			SwiftCode:     "ABCDEF12XXX",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: true,
		}

		// when
		err := branchRepository.CreateBranch(branch)

		// then
		assert.NoError(t, err)

		var storedBranch model.Branch
		err = db.First(&storedBranch, "swift_code = ?", branch.SwiftCode).Error
		assert.NoError(t, err)
		assert.Equal(t, branch, &storedBranch)
	})

	t.Run("should return error when creating branch with duplicate SwiftCode", func(t *testing.T) {
		// given
		duplicateBranch := &model.Branch{
			SwiftCode:     "ABCDEF12XXX",
			Address:       "Another Address",
			BankName:      "Another Bank",
			CountryISO2:   "DE",
			CountryName:   "Germany",
			IsHeadquarter: true,
		}

		// when
		err := branchRepository.CreateBranch(duplicateBranch)

		// then
		assert.Error(t, err)
	})
}

func TestBranchRepository_GetBranchBySwiftCode(t *testing.T) {
	db := setupTestDB()
	branchRepository := NewBranchRepository(db)

	t.Run("should get branch successfully", func(t *testing.T) {
		// given
		branch1 := &model.Branch{
			SwiftCode:     "AIPOPLP1XXX",
			Address:       "123 Test St",
			BankName:      "Test Bank1",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: true,
		}
		swiftCode1 := branch1.SwiftCode

		branch2 := &model.Branch{
			SwiftCode:     "AIPOPLP1FGD",
			Address:       "456 Another St",
			BankName:      "Another Bank",
			CountryISO2:   "PL",
			CountryName:   "Germany",
			IsHeadquarter: false,
		}
		swiftCode2 := branch2.SwiftCode

		// when
		resultBranch1, err1 := branchRepository.GetBranchBySwiftCode(swiftCode1)
		resultBranch2, err2 := branchRepository.GetBranchBySwiftCode(swiftCode2)

		// then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, branch1, resultBranch1)
		assert.Equal(t, branch2, resultBranch2)
	})

	t.Run("should return nil when branch with given swift-code does not exist", func(t *testing.T) {
		// given
		swiftCode1 := "BIPOPLP1XXX"
		swiftCode2 := "BIPOPLP1ABC"

		// when
		resultBranch1, err1 := branchRepository.GetBranchBySwiftCode(swiftCode1)
		resultBranch2, err2 := branchRepository.GetBranchBySwiftCode(swiftCode2)

		// then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Nil(t, resultBranch1)
		assert.Nil(t, resultBranch2)

	})

}

func TestBranchRepository_GetBranchesByISO2code(t *testing.T) {
	db := setupTestDB()
	branchRepository := NewBranchRepository(db)

	t.Run("should get two branches successfully", func(t *testing.T) {
		// given
		isoCodeWithTwoBranches := "PL"

		twoBranches := []model.Branch{
			{
				SwiftCode:     "AIPOPLP1XXX",
				Address:       "123 Test St",
				BankName:      "Test Bank1",
				CountryISO2:   "PL",
				CountryName:   "Poland",
				IsHeadquarter: true,
			},
			{
				SwiftCode:     "AIPOPLP1FGD",
				Address:       "456 Another St",
				BankName:      "Another Bank",
				CountryISO2:   "PL",
				CountryName:   "Germany",
				IsHeadquarter: false,
			},
		}

		// when
		branches, err := branchRepository.GetBranchesByISO2code(isoCodeWithTwoBranches)

		// then
		assert.NoError(t, err)
		assert.Equal(t, twoBranches, branches)

	})

	t.Run("should return one branch in a list successfully", func(t *testing.T) {
		// given
		isoCodeWithOneBranch := "BG"
		oneBranch := []model.Branch{
			{
				SwiftCode:     "FEDCBA12XXX",
				Address:       "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",
				BankName:      "UNITED BANK OF ALBANIA SH.A",
				CountryISO2:   "BG",
				CountryName:   "Germany",
				IsHeadquarter: true,
			},
		}

		// when
		branches, err := branchRepository.GetBranchesByISO2code(isoCodeWithOneBranch)

		// then
		assert.NoError(t, err)
		assert.Equal(t, oneBranch, branches)
	})

	t.Run("should return empty list because there is no branches with such ISO2 code", func(t *testing.T) {
		// given
		isoCodeWithoutBranches := "AB"

		// when
		branches, err := branchRepository.GetBranchesByISO2code(isoCodeWithoutBranches)

		// then
		assert.NoError(t, err)
		assert.Equal(t, []model.Branch{}, branches)
	})
}

func TestBranchRepository_GetOrdinaryBranchesForHeadquarter(t *testing.T) {
	db := setupTestDB()
	branchRepository := NewBranchRepository(db)

	t.Run("should return branches properly for headquarter", func(t *testing.T) {
		// given
		headquarterSwift := "AIPOPLP1XXX"
		expectedBranches := []model.Branch{
			{
				SwiftCode:     "AIPOPLP1FGD",
				Address:       "456 Another St",
				BankName:      "Another Bank",
				CountryISO2:   "PL",
				CountryName:   "Germany",
				IsHeadquarter: false,
			},
		}

		// when
		resultBranches, err := branchRepository.GetOrdinaryBranchesForHeadquarter(headquarterSwift)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedBranches, resultBranches)
	})

	t.Run("should return empty list because headquarter does not have branches", func(t *testing.T) {
		// given
		headquarterSwift := "FEDCBA12XXX"

		// when
		resultBraches, err := branchRepository.GetOrdinaryBranchesForHeadquarter(headquarterSwift)

		// then
		assert.NoError(t, err)
		assert.Equal(t, []model.Branch{}, resultBraches)
	})

}

func TestBranchRepository_RemoveBranchBySwiftCode(t *testing.T) {
	db := setupTestDB()
	branchRepository := NewBranchRepository(db)

	t.Run("should remove branch successfully", func(t *testing.T) {
		// given
		branch1 := &model.Branch{
			SwiftCode:     "AIPOPLP1XXX",
			Address:       "123 Test St",
			BankName:      "Test Bank1",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: true,
		}

		branch2 := &model.Branch{
			SwiftCode:     "AIPOPLP1FGD",
			Address:       "456 Another St",
			BankName:      "Another Bank",
			CountryISO2:   "PL",
			CountryName:   "Germany",
			IsHeadquarter: false,
		}

		// when & then
		resultBranch, _ := branchRepository.GetBranchBySwiftCode(branch1.SwiftCode)
		assert.Equal(t, branch1, resultBranch)
		err := branchRepository.RemoveBranchBySwiftCode(branch1.SwiftCode)
		assert.NoError(t, err)
		resultBranch, _ = branchRepository.GetBranchBySwiftCode(branch1.SwiftCode)
		assert.Nil(t, resultBranch)

		resultBranch2, _ := branchRepository.GetBranchBySwiftCode(branch2.SwiftCode)
		assert.Equal(t, branch2, resultBranch2)
		err = branchRepository.RemoveBranchBySwiftCode(branch2.SwiftCode)
		assert.NoError(t, err)
		resultBranch2, _ = branchRepository.GetBranchBySwiftCode(branch2.SwiftCode)
		assert.Nil(t, resultBranch2)
	})

}
