package handler

import (
	"SWIFT_task/internal/db"
	"SWIFT_task/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetBranchDetails(c *gin.Context) {
	swiftCode := c.Param("swift-code")
	var branch model.Branch

	if err := db.DB.Where("swift_code = ?", swiftCode).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	branchDto := branch.ToBranchDto()

	if branch.IsHeadquarter {
		var branches []model.Branch

		db.DB.Joins("JOIN branch_relationships ON branches.swift_code = branch_relationships.ordinary_branch_swift_code").
			Where("branch_relationships.headquarter_swift_code = ?", swiftCode).
			Find(&branches)

		branchDTOs := make([]model.BranchWithoutCountryNameDto, len(branches))
		for i, branch := range branches {
			branchDTOs[i] = branch.ToBranchWithoutCountryNameDto()
		}

		branchDto.Branches = branchDTOs

		c.JSON(http.StatusOK, branchDto)
		return
	}

	c.JSON(http.StatusOK, branchDto)

}

func GetBranchesByISO2code(c *gin.Context) {
	countryISO2code := c.Param("countryISO2code")

	var branches []model.Branch
	if err := db.DB.Where("country_iso2 = ?", countryISO2code).Find(&branches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if len(branches) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No branches found for the provided ISO code"})
		return
	}

	countryName := branches[0].CountryName
	branchesDto := make([]model.BranchWithoutCountryNameDto, len(branches))
	for i, branch := range branches {
		branchesDto[i] = branch.ToBranchWithoutCountryNameDto()
	}

	responseDto := model.BranchesForCountryDto{
		CountryISO2: countryISO2code,
		CountryName: countryName,
		Branches:    branchesDto,
	}

	c.JSON(http.StatusOK, responseDto)
}
