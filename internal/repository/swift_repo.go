package repository

import (
	"SWIFT_task/internal/model"
	"errors"
	"gorm.io/gorm"
)

type BranchRepository interface {
	GetBranchBySwiftCode(swiftCode string) (*model.Branch, error)
	GetOrdinaryBranchesForHeadquarter(swiftCode string) ([]model.Branch, error)
	GetBranchesByISO2code(countryISO2code string) ([]model.Branch, error)
	CreateBranch(branch *model.Branch) error
	RemoveBranchBySwiftCode(swiftCode string) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) GetBranchBySwiftCode(swiftCode string) (*model.Branch, error) {
	var branch model.Branch
	if err := r.db.Where("swift_code = ?", swiftCode).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) GetOrdinaryBranchesForHeadquarter(swiftCode string) ([]model.Branch, error) {
	var branches []model.Branch
	err := r.db.Joins("JOIN branch_relationships ON branches.swift_code = branch_relationships.ordinary_branch_swift_code").
		Where("branch_relationships.headquarter_swift_code = ?", swiftCode).
		Find(&branches).Error

	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *branchRepository) GetBranchesByISO2code(countryISO2code string) ([]model.Branch, error) {
	var branches []model.Branch
	err := r.db.Where("country_iso2 = ?", countryISO2code).
		Find(&branches).
		Error

	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *branchRepository) CreateBranch(branch *model.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) RemoveBranchBySwiftCode(swiftCode string) error {
	return r.db.Where("swift_code = ?", swiftCode).Delete(&model.Branch{}).Error
}
