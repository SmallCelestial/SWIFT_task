package handler

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type BranchHandler struct {
	branchService *service.BranchService
}

func NewBranchHandler(branchService *service.BranchService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

func (h *BranchHandler) GetBranchDetails(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	branchDto, err := h.branchService.GetBranchDetails(swiftCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if branchDto == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	c.JSON(http.StatusOK, branchDto)
}

func (h *BranchHandler) GetBranchesByISO2code(c *gin.Context) {
	countryISO2code := c.Param("countryISO2code")

	response, err := h.branchService.GetBranchesByISO2code(countryISO2code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branches not found for ISO2 code " + countryISO2code})
	}

	c.JSON(http.StatusOK, response)
}

func (h *BranchHandler) AddSwiftCode(c *gin.Context) {

	var branch model.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	if err := h.branchService.AddSwiftCode(&branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add SWIFT code", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code added successfully"})
}

func (h *BranchHandler) RemoveSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	if err := h.branchService.RemoveBranchBySwiftCode(swiftCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code " + swiftCode + " not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove SWIFT code", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code " + swiftCode + " removed successfully"})
}
