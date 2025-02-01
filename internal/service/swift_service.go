package service

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/repository"
	"fmt"
	"strconv"
)

type BranchService struct {
	branchRepo repository.BranchRepository
}

func NewBranchService(branchRepo repository.BranchRepository) *BranchService {
	return &BranchService{branchRepo: branchRepo}
}

func (s *BranchService) GetBranchDetails(swiftCode string) (*model.BankDto, error) {
	branch, err := s.branchRepo.GetBranchBySwiftCode(swiftCode)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, NewBranchNotExistsError("Bank with swiftCode " + swiftCode + " not exists.")
	}

	branchDto := branch.ToBranchDto()

	if branch.IsHeadquarter {
		branches, err := s.branchRepo.GetOrdinaryBranchesForHeadquarter(swiftCode)
		if err != nil {
			return nil, err
		}

		branchDTOs := make([]model.BankWithoutCountryNameDto, len(branches))
		for i, b := range branches {
			branchDTOs[i] = b.ToBranchWithoutCountryNameDto()
		}

		branchDto.Branches = branchDTOs
	}

	return &branchDto, nil
}

func (s *BranchService) GetBranchesByISO2code(countryISO2code string) (*model.CountryBanksDto, error) {
	if len(countryISO2code) != 2 {
		return nil, NewValidationError("Len of countryISO2 code should be 2. Got len: " + strconv.Itoa(len(countryISO2code)))
	}

	branches, err := s.branchRepo.GetBranchesByISO2code(countryISO2code)
	if err != nil {
		return nil, err
	}
	if len(branches) == 0 {
		return nil, nil
	}

	countryName := branches[0].CountryName
	branchesDto := make([]model.BankWithoutCountryNameDto, len(branches))
	for i, branch := range branches {
		branchesDto[i] = branch.ToBranchWithoutCountryNameDto()
	}

	return &model.CountryBanksDto{
		CountryISO2: countryISO2code,
		CountryName: countryName,
		Branches:    branchesDto,
	}, err
}

func (s *BranchService) AddSwiftCode(branch model.Bank) error {
	if len(branch.CountryISO2) != 2 {
		return NewValidationError("Len of countryISO2 code should be 2. Got len: " + strconv.Itoa(len(branch.CountryISO2)))
	}

	existingBranch, err := s.branchRepo.GetBranchBySwiftCode(branch.SwiftCode)
	if err != nil {
		return err
	}
	if existingBranch != nil {
		return NewBranchExistsError("Bank with swiftCode " + branch.SwiftCode + " already exists.")
	}
	fmt.Println("Adding swift code: " + branch.SwiftCode)
	return s.branchRepo.CreateBranch(&branch)
}

func (s *BranchService) RemoveBranchBySwiftCode(swiftCode string) error {
	branch, err := s.branchRepo.GetBranchBySwiftCode(swiftCode)
	if err != nil {
		return err
	}
	if branch == nil {
		return NewBranchNotExistsError("Bank with swiftCode " + swiftCode + " not exists.")
	}

	return s.branchRepo.RemoveBranchBySwiftCode(swiftCode)
}
