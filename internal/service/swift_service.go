package service

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/repository"
)

type BranchService struct {
	branchRepo *repository.BranchRepository
}

func NewBranchService(branchRepo *repository.BranchRepository) *BranchService {
	return &BranchService{branchRepo: branchRepo}
}

func (s *BranchService) GetBranchDetails(swiftCode string) (*model.BranchDto, error) {
	branch, err := s.branchRepo.GetBranchBySwiftCode(swiftCode)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, nil
	}

	branchDto := branch.ToBranchDto()

	if branch.IsHeadquarter {
		branches, err := s.branchRepo.GetOrdinaryBranchesForHeadquarter(swiftCode)
		if err != nil {
			return nil, err
		}

		branchDTOs := make([]model.BranchWithoutCountryNameDto, len(branches))
		for i, b := range branches {
			branchDTOs[i] = b.ToBranchWithoutCountryNameDto()
		}

		branchDto.Branches = branchDTOs
	}

	return &branchDto, nil
}

func (s *BranchService) GetBranchesByISO2code(countryISO2code string) (*model.BranchesForCountryDto, error) {
	branches, err := s.branchRepo.GetBranchesByISO2code(countryISO2code)
	if err != nil {
		return nil, err
	}
	if len(branches) == 0 {
		return nil, nil
	}

	countryName := branches[0].CountryName
	branchesDto := make([]model.BranchWithoutCountryNameDto, len(branches))
	for i, branch := range branches {
		branchesDto[i] = branch.ToBranchWithoutCountryNameDto()
	}

	return &model.BranchesForCountryDto{
		CountryISO2: countryISO2code,
		CountryName: countryName,
		Branches:    branchesDto,
	}, err
}

func (s *BranchService) AddSwiftCode(branch *model.Branch) error {
	return s.branchRepo.CreateBranch(branch)
}

func (s *BranchService) RemoveBranchBySwiftCode(swiftCode string) error {
	_, err := s.branchRepo.GetBranchBySwiftCode(swiftCode)
	if err != nil {
		return err
	}
	return s.branchRepo.RemoveBranchBySwiftCode(swiftCode)
}
