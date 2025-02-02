package handler

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BranchHandler struct {
	branchService *service.BankService
}

func NewBranchHandler(branchService *service.BankService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

func (h *BranchHandler) GetBranchDetails(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	branchDto, err := h.branchService.GetBankDetails(swiftCode)
	var branchServiceErr *service.BankNotExistsError
	if errors.As(err, &branchServiceErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
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

	response, err := h.branchService.GetBanksByISO2code(countryISO2code)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Banks not found for ISO2 code " + countryISO2code})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *BranchHandler) AddSwiftCode(c *gin.Context) {

	var branch model.Bank
	err := c.ShouldBindJSON(&branch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	err = h.branchService.AddBank(branch)

	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var branchExistsErr *service.BankExistsError
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

	err := h.branchService.RemoveBankBySwiftCode(swiftCode)
	var branchNotExistsErr *service.BankNotExistsError
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
