package handler

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
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
	var branchServiceErr *service.BranchNotExistsError
	if errors.As(err, &branchServiceErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, branchDto)
}

func (h *BranchHandler) GetBranchesByISO2code(c *gin.Context) {
	countryISO2code := c.Param("countryISO2code")

	response, err := h.branchService.GetBranchesByISO2code(countryISO2code)
	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branches not found for ISO2 code " + countryISO2code})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *BranchHandler) AddSwiftCode(c *gin.Context) {

	var branch model.Branch
	err := c.ShouldBindJSON(&branch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	err = h.branchService.AddSwiftCode(branch)

	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var branchExistsErr *service.BranchExistsError
	if errors.As(err, &branchExistsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": branchExistsErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add SWIFT code", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code added successfully"})
}

func (h *BranchHandler) RemoveSwiftCode(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	err := h.branchService.RemoveBranchBySwiftCode(swiftCode)
	var branchNotExistsErr *service.BranchNotExistsError
	if errors.As(err, &branchNotExistsErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code " + swiftCode + " not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove SWIFT code", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code " + swiftCode + " removed successfully"})
}
