package service

import "fmt"

type ValidationError struct {
	Message string
}

type BankExistsError struct {
	Message string
}

type BankNotExistsError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func (e *BankExistsError) Error() string {
	return fmt.Sprintf("bank exists: %s", e.Message)
}

func NewBankExistsError(message string) *BankExistsError {
	return &BankExistsError{Message: message}
}

func (e *BankNotExistsError) Error() string {
	return fmt.Sprintf("bank not found: %s", e.Message)
}

func NewBankNotExistsError(message string) *BankNotExistsError {
	return &BankNotExistsError{Message: message}
}
