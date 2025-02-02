package service

import (
	"SWIFT_task/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockBranchRepo struct {
	mock.Mock
}

func (m *MockBranchRepo) GetBankBySwiftCode(swiftCode string) (*model.Bank, error) {
	args := m.Called(swiftCode)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Bank), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBranchRepo) GetBranchesForHeadquarter(swiftCode string) ([]model.Bank, error) {
	args := m.Called(swiftCode)
	return args.Get(0).([]model.Bank), args.Error(1)
}

func (m *MockBranchRepo) GetBanksByISO2code(countryISO2code string) ([]model.Bank, error) {
	args := m.Called(countryISO2code)
	return args.Get(0).([]model.Bank), args.Error(1)
}

func (m *MockBranchRepo) AddBank(branch *model.Bank) error {
	args := m.Called(branch)
	return args.Error(0)
}

func (m *MockBranchRepo) RemoveBankBySwiftCode(swiftCode string) error {
	args := m.Called(swiftCode)
	return args.Error(0)
}

func TestBankService_GetBankDetails(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	service := &BankService{bankRepo: mockRepo}

	t.Run("Should return BankNotExistError due to not existing bank with given swiftCode", func(t *testing.T) {
		// given
		mockRepo.On("GetBankBySwiftCode", "XXXXXX12346").Return(nil, nil)

		// when
		branch, err := service.GetBankDetails("XXXXXX12346")

		// then
		var bankNotExistsErr *BankNotExistsError
		assert.ErrorAs(t, err, &bankNotExistsErr)
		assert.Equal(t, "Bank with swiftCode XXXXXX12346 not exists.", bankNotExistsErr.Message)
		assert.Nil(t, branch)
	})

	t.Run("should return branchInfo", func(t *testing.T) {
		// given
		validBank := model.Bank{
			SwiftCode:   "ABCDEF12346",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		expectedBranchInfo := validBank.ToBankDto()
		mockRepo.On("GetBankBySwiftCode", validBank.SwiftCode).Return(&validBank, nil)

		// when
		branch, err := service.GetBankDetails(validBank.SwiftCode)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedBranchInfo, branch)
	})

	t.Run("should return headquarterInfo", func(t *testing.T) {
		// given
		headquarter := model.Bank{
			SwiftCode:   "AAAAA12XXX",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		branch1 := &model.Bank{
			SwiftCode:   "AAAAA12346",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		branch2 := &model.Bank{
			SwiftCode:   "AAAAA12RTP",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		mockRepo.On("GetBankBySwiftCode", headquarter.SwiftCode).Return(&headquarter, nil)
		mockRepo.On("GetBankBySwiftCode", branch1.SwiftCode).Return(&branch1, nil)
		mockRepo.On("GetBankBySwiftCode", branch2.SwiftCode).Return(&branch2, nil)
		mockRepo.On("GetBranchesForHeadquarter", headquarter.SwiftCode).Return([]model.Bank{*branch1, *branch2}, nil)
		expectedBankInfo := headquarter.ToBankDto()
		expectedBankInfo.Branches = []model.BankWithoutCountryNameDto{branch1.ToBankWithoutCountryNameDto(), branch2.ToBankWithoutCountryNameDto()}

		// when
		result, err := service.GetBankDetails(headquarter.SwiftCode)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedBankInfo, result)

	})

}

func TestBankService_GetBanksByISO2code(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	service := &BankService{bankRepo: mockRepo}

	t.Run("Should return ValidationError due to ISO2 code", func(t *testing.T) {
		// given
		countryISO2 := "PL34"

		// when
		resultBanks, err := service.GetBanksByISO2code(countryISO2)

		// then
		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "Len of countryISO2 code should be 2. Got len: 4", validationErr.Message)
		assert.Nil(t, resultBanks)
	})

	t.Run("Should return Banks with given ISO2 code in proper structure", func(t *testing.T) {
		// given
		bank1 := model.Bank{
			SwiftCode:   "AAAAA12XXX",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		bank2 := model.Bank{
			SwiftCode:   "AAAAA12346",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		mockRepo.On("GetBanksByISO2code", "PL").Return([]model.Bank{bank1, bank2}, nil)
		expectedBanks := model.CountryBanksDto{
			CountryISO2: "PL",
			CountryName: "POLAND",
			Banks:       []model.BankWithoutCountryNameDto{bank1.ToBankWithoutCountryNameDto(), bank2.ToBankWithoutCountryNameDto()},
		}

		// when
		banks, err := service.GetBanksByISO2code("PL")

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedBanks, banks)

	})

}

func TestBankService_AddBank(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	mockRepo.On("AddBank", mock.Anything).Return(nil)
	service := &BankService{bankRepo: mockRepo}

	t.Run("Should return ValidationError due to ISO2 code", func(t *testing.T) {
		// given
		bankWithTooLongCountryISO2 := model.Bank{
			SwiftCode:   "ABCDEF12345",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL3",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}

		// when
		err := service.AddBank(bankWithTooLongCountryISO2)

		// then
		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "Len of countryISO2 code should be 2. Got len: 3", validationErr.Message)
	})

	t.Run("Should return BankExistsError due to duplicate key", func(t *testing.T) {
		// given
		duplicatedBank := model.Bank{
			SwiftCode:   "ABCDEF12345",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		mockRepo.On("GetBankBySwiftCode", duplicatedBank.SwiftCode).Return(&duplicatedBank, nil)

		// when
		err := service.AddBank(duplicatedBank)

		// then
		var bankExistsErr *BankExistsError
		assert.ErrorAs(t, err, &bankExistsErr)
		assert.Equal(t, "Bank with swiftCode "+duplicatedBank.SwiftCode+" already exists.", bankExistsErr.Message)
	})

	t.Run("Should properly add bank", func(t *testing.T) {
		// given
		validBank := model.Bank{
			SwiftCode:   "ABCDEF12346",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		mockRepo.On("GetBankBySwiftCode", validBank.SwiftCode).Return(nil, nil)

		// when
		err := service.AddBank(validBank)

		//then
		assert.NoError(t, err)
	})

}

func TestBankService_RemoveBranchBySwiftCode(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	service := &BankService{bankRepo: mockRepo}

	t.Run("Should return BankNotExistError due to bank doesn't exist", func(t *testing.T) {
		// given
		mockRepo.On("GetBankBySwiftCode", "AAAAAAAAAAA").Return(nil, nil)
		// when
		err := service.RemoveBankBySwiftCode("AAAAAAAAAAA")

		// then
		var bankNotExistsError *BankNotExistsError
		assert.ErrorAs(t, err, &bankNotExistsError)
		assert.Equal(t, "Bank with swiftCode AAAAAAAAAAA not exists.", bankNotExistsError.Message)
	})

	t.Run("Should remove bank successfully ", func(t *testing.T) {
		// given
		bank := model.Bank{
			SwiftCode:   "BBBBBBBBBBB",
			Address:     "123 Test St",
			BankName:    "Test Bank",
			CountryISO2: "PL",
			Country: model.Country{
				CountryISO2: "PL",
				CountryName: "POLAND",
			},
		}
		mockRepo.On("GetBankBySwiftCode", "BBBBBBBBBBB").Return(&bank, nil)
		mockRepo.On("RemoveBankBySwiftCode", "BBBBBBBBBBB").Return(nil)

		// when
		err := service.RemoveBankBySwiftCode("BBBBBBBBBBB")

		// then
		assert.NoError(t, err)

	})

}
