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
