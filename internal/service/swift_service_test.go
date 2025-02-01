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

func (m *MockBranchRepo) GetBranchBySwiftCode(swiftCode string) (*model.Bank, error) {
	args := m.Called(swiftCode)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Bank), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBranchRepo) GetOrdinaryBranchesForHeadquarter(swiftCode string) ([]model.Bank, error) {
	args := m.Called(swiftCode)
	return args.Get(0).([]model.Bank), args.Error(1)
}

func (m *MockBranchRepo) GetBranchesByISO2code(countryISO2code string) ([]model.Bank, error) {
	args := m.Called(countryISO2code)
	return args.Get(0).([]model.Bank), args.Error(1)
}

func (m *MockBranchRepo) CreateBranch(branch *model.Bank) error {
	args := m.Called(branch)
	return args.Error(0)
}

func (m *MockBranchRepo) RemoveBranchBySwiftCode(swiftCode string) error {
	args := m.Called(swiftCode)
	return args.Error(0)
}

func TestBranchService_AddSwiftCode(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	mockRepo.On("CreateBranch", mock.Anything).Return(nil)
	service := &BranchService{branchRepo: mockRepo}

	t.Run("Should return ValidationError due to ISO2 code", func(t *testing.T) {
		// given
		branchWithTooLongCountryISO2 := model.Bank{
			SwiftCode:     "ABCDEF12345",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL3",
			CountryName:   "Poland",
			IsHeadquarter: true,
		}

		// when
		err := service.AddSwiftCode(branchWithTooLongCountryISO2)

		// then
		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "Len of countryISO2 code should be 2. Got len: 3", validationErr.Message)
	})

	t.Run("Should return BranchExistsError due to duplicate key", func(t *testing.T) {
		// given
		duplicatedBranch := model.Bank{
			SwiftCode:     "ABCDEF12345",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		mockRepo.On("GetBranchBySwiftCode", duplicatedBranch.SwiftCode).Return(&duplicatedBranch, nil)

		// when
		err := service.AddSwiftCode(duplicatedBranch)

		// then
		var branchExistsErr *BranchExistsError
		assert.ErrorAs(t, err, &branchExistsErr)
		assert.Equal(t, "Bank with swiftCode "+duplicatedBranch.SwiftCode+" already exists.", branchExistsErr.Message)
	})

	t.Run("Should properly add branch", func(t *testing.T) {
		// given
		validBranch := model.Bank{
			SwiftCode:     "ABCDEF12346",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		mockRepo.On("GetBranchBySwiftCode", validBranch.SwiftCode).Return(nil, nil)

		// when
		err := service.AddSwiftCode(validBranch)

		//then
		assert.NoError(t, err)
	})

}

func TestBranchService_GetBranchDetails(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	service := &BranchService{branchRepo: mockRepo}

	t.Run("Should return BranchNotExistError due to not existing branch with given swiftCode", func(t *testing.T) {
		// given
		mockRepo.On("GetBranchBySwiftCode", "XXXXXX12346").Return(nil, nil)

		// when
		branch, err := service.GetBranchDetails("XXXXXX12346")

		// then
		var branchNotExistsErr *BranchNotExistsError
		assert.ErrorAs(t, err, &branchNotExistsErr)
		assert.Equal(t, "Bank with swiftCode XXXXXX12346 not exists.", branchNotExistsErr.Message)
		assert.Nil(t, branch)
	})

	t.Run("should return ordinary branchInfo", func(t *testing.T) {
		// given
		validBranch := model.Bank{
			SwiftCode:     "ABCDEF12346",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		expectedBranchInfo := validBranch.ToBranchDto()
		mockRepo.On("GetBranchBySwiftCode", validBranch.SwiftCode).Return(&validBranch, nil)

		// when
		branch, err := service.GetBranchDetails(validBranch.SwiftCode)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedBranchInfo, branch)
	})

	t.Run("should return headquarterInfo", func(t *testing.T) {
		// given
		headquarter := model.Bank{
			SwiftCode:     "AAAAA12XXX",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: true,
		}
		branch1 := &model.Bank{
			SwiftCode:     "AAAAA12346",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		branch2 := &model.Bank{
			SwiftCode:     "AAAAA12RTP",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		mockRepo.On("GetBranchBySwiftCode", headquarter.SwiftCode).Return(&headquarter, nil)
		mockRepo.On("GetBranchBySwiftCode", branch1.SwiftCode).Return(&branch1, nil)
		mockRepo.On("GetBranchBySwiftCode", branch2.SwiftCode).Return(&branch2, nil)
		mockRepo.On("GetOrdinaryBranchesForHeadquarter", headquarter.SwiftCode).Return([]model.Bank{*branch1, *branch2}, nil)
		expectedBranchInfo := headquarter.ToBranchDto()
		expectedBranchInfo.Branches = []model.BankWithoutCountryNameDto{branch1.ToBranchWithoutCountryNameDto(), branch2.ToBranchWithoutCountryNameDto()}

		// when
		result, err := service.GetBranchDetails(headquarter.SwiftCode)

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedBranchInfo, result)

	})

}

func TestBranchService_GetBranchesByISO2code(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	service := &BranchService{branchRepo: mockRepo}

	t.Run("Should return ValidationError due to ISO2 code", func(t *testing.T) {
		// given
		countryISO2 := "PL34"

		// when
		resultBranches, err := service.GetBranchesByISO2code(countryISO2)

		// then
		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "Len of countryISO2 code should be 2. Got len: 4", validationErr.Message)
		assert.Nil(t, resultBranches)
	})

	t.Run("Should return Branches with given ISO2 code in proper structure", func(t *testing.T) {
		// given
		branch1 := model.Bank{
			SwiftCode:     "AAAAA12XXX",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: true,
		}
		branch2 := model.Bank{
			SwiftCode:     "AAAAA12346",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		mockRepo.On("GetBranchesByISO2code", "PL").Return([]model.Bank{branch1, branch2}, nil)
		expectedBranches := model.CountryBanksDto{
			CountryISO2: "PL",
			CountryName: "Poland",
			Branches:    []model.BankWithoutCountryNameDto{branch1.ToBranchWithoutCountryNameDto(), branch2.ToBranchWithoutCountryNameDto()},
		}

		// when
		branches, err := service.GetBranchesByISO2code("PL")

		// then
		assert.NoError(t, err)
		assert.Equal(t, &expectedBranches, branches)

	})

}

func TestBranchService_RemoveBranchBySwiftCode(t *testing.T) {
	mockRepo := new(MockBranchRepo)
	service := &BranchService{branchRepo: mockRepo}

	t.Run("Should return BranchNotExistError due to branch doesn't exist", func(t *testing.T) {
		// given
		mockRepo.On("GetBranchBySwiftCode", "AAAAAAAAAAA").Return(nil, nil)
		// when
		err := service.RemoveBranchBySwiftCode("AAAAAAAAAAA")

		// then
		var branchNotExistErr *BranchNotExistsError
		assert.ErrorAs(t, err, &branchNotExistErr)
		assert.Equal(t, "Bank with swiftCode AAAAAAAAAAA not exists.", branchNotExistErr.Message)
	})

	t.Run("Should remove branch successfully ", func(t *testing.T) {
		// given
		branch := model.Bank{
			SwiftCode:     "BBBBBBBBBBB",
			Address:       "123 Test St",
			BankName:      "Test Bank",
			CountryISO2:   "PL",
			CountryName:   "Poland",
			IsHeadquarter: false,
		}
		mockRepo.On("GetBranchBySwiftCode", "BBBBBBBBBBB").Return(&branch, nil)
		mockRepo.On("RemoveBranchBySwiftCode", "BBBBBBBBBBB").Return(nil)

		// when
		err := service.RemoveBranchBySwiftCode("BBBBBBBBBBB")

		// then
		assert.NoError(t, err)

	})

}
