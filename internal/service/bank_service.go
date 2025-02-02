package service

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/repository"
	"strconv"
)

type BankService struct {
	bankRepo repository.BankRepository
}

func NewBankService(bankRepo repository.BankRepository) *BankService {
	return &BankService{bankRepo: bankRepo}
}

func (s *BankService) GetBankDetails(swiftCode string) (*model.BankDto, error) {
	bank, err := s.bankRepo.GetBankBySwiftCode(swiftCode)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, NewBankNotExistsError("Bank with swiftCode " + swiftCode + " not exists.")
	}

	bankDto := bank.ToBankDto()

	if bank.IsHeadquarter() {
		branches, err := s.bankRepo.GetBranchesForHeadquarter(swiftCode)
		if err != nil {
			return nil, err
		}

		branchDTOs := make([]model.BankWithoutCountryNameDto, len(branches))
		for i, b := range branches {
			branchDTOs[i] = b.ToBankWithoutCountryNameDto()
		}

		bankDto.Branches = branchDTOs
	}

	return &bankDto, nil
}

func (s *BankService) GetBanksByISO2code(countryISO2code string) (*model.CountryBanksDto, error) {
	if len(countryISO2code) != 2 {
		return nil, NewValidationError("Len of countryISO2 code should be 2. Got len: " + strconv.Itoa(len(countryISO2code)))
	}

	banks, err := s.bankRepo.GetBanksByISO2code(countryISO2code)
	if err != nil {
		return nil, err
	}
	if len(banks) == 0 {
		return nil, nil
	}

	countryName := banks[0].Country.CountryName
	banksDto := make([]model.BankWithoutCountryNameDto, len(banks))
	for i, branch := range banks {
		banksDto[i] = branch.ToBankWithoutCountryNameDto()
	}

	return &model.CountryBanksDto{
		CountryISO2: countryISO2code,
		CountryName: countryName,
		Banks:       banksDto,
	}, err
}

func (s *BankService) AddBank(bank model.Bank) error {
	if len(bank.CountryISO2) != 2 {
		return NewValidationError("Len of countryISO2 code should be 2. Got len: " + strconv.Itoa(len(bank.CountryISO2)))
	}

	existingBank, err := s.bankRepo.GetBankBySwiftCode(bank.SwiftCode)
	if err != nil {
		return err
	}
	if existingBank != nil {
		return NewBankExistsError("Bank with swiftCode " + bank.SwiftCode + " already exists.")
	}
	return s.bankRepo.AddBank(&bank)
}

func (s *BankService) RemoveBankBySwiftCode(swiftCode string) error {
	bank, err := s.bankRepo.GetBankBySwiftCode(swiftCode)
	if err != nil {
		return err
	}
	if bank == nil {
		return NewBankNotExistsError("Bank with swiftCode " + swiftCode + " not exists.")
	}

	return s.bankRepo.RemoveBankBySwiftCode(swiftCode)
}
