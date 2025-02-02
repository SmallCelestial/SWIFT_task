package handler

import (
	"SWIFT_task/internal/model"
	"SWIFT_task/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BankHandler struct {
	bankService *service.BankService
}

func NewBankHandler(bankService *service.BankService) *BankHandler {
	return &BankHandler{bankService: bankService}
}

func (h *BankHandler) GetBankDetails(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	bankDto, err := h.bankService.GetBankDetails(swiftCode)
	var bankServiceErr *service.BankNotExistsError
	if errors.As(err, &bankServiceErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, bankDto)
}

func (h *BankHandler) GetBanksByISO2code(c *gin.Context) {
	countryISO2code := c.Param("countryISO2code")

	response, err := h.bankService.GetBanksByISO2code(countryISO2code)
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

func (h *BankHandler) AddBank(c *gin.Context) {

	var bank model.Bank
	err := c.ShouldBindJSON(&bank)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	err = h.bankService.AddBank(bank)

	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var bankExistsErr *service.BankExistsError
	if errors.As(err, &bankExistsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": bankExistsErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add new bank", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank added successfully"})
}

func (h *BankHandler) RemoveBank(c *gin.Context) {
	swiftCode := c.Param("swift-code")

	err := h.bankService.RemoveBankBySwiftCode(swiftCode)
	var bankNotExistsErr *service.BankNotExistsError
	if errors.As(err, &bankNotExistsErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bank with swift code: " + swiftCode + " not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove bank", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank with swift code:  " + swiftCode + " removed successfully"})
}
