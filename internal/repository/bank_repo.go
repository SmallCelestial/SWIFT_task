package repository

import (
	"SWIFT_task/internal/model"
	"errors"
	"gorm.io/gorm"
)

type BankRepository interface {
	GetBankBySwiftCode(swiftCode string) (*model.Bank, error)
	GetBranchesForHeadquarter(swiftCode string) ([]model.Bank, error)
	GetBanksByISO2code(countryISO2code string) ([]model.Bank, error)
	AddBank(branch *model.Bank) error
	RemoveBankBySwiftCode(swiftCode string) error
}

type bankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &bankRepository{db: db}
}

func (r *bankRepository) GetBankBySwiftCode(swiftCode string) (*model.Bank, error) {
	var bank model.Bank
	err := r.db.
		Preload("Country").
		Where("swift_code = ?", swiftCode).
		First(&bank).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (r *bankRepository) GetBranchesForHeadquarter(swiftCode string) ([]model.Bank, error) {
	var banks []model.Bank
	err := r.db.
		Preload("Country").
		Joins("JOIN bank_relationships br ON banks.swift_code = br.branch_swift_code").
		Where("br.headquarter_swift_code = ?", swiftCode).
		Find(&banks).Error

	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (r *bankRepository) GetBanksByISO2code(countryISO2code string) ([]model.Bank, error) {
	var banks []model.Bank
	err := r.db.
		Preload("Country").
		Where("country_iso2 = ?", countryISO2code).
		Find(&banks).
		Error

	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (r *bankRepository) AddBank(bank *model.Bank) error {
	return r.db.Create(bank).Error
}

func (r *bankRepository) RemoveBankBySwiftCode(swiftCode string) error {
	return r.db.Where("swift_code = ?", swiftCode).Delete(&model.Bank{}).Error
}
