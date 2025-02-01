package service

import "fmt"

type ValidationError struct {
	Message string
}

type BranchExistsError struct {
	Message string
}

type BranchNotExistsError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func (e *BranchExistsError) Error() string {
	return fmt.Sprintf("branch exists: %s", e.Message)
}

func NewBranchExistsError(message string) *BranchExistsError {
	return &BranchExistsError{Message: message}
}

func (e *BranchNotExistsError) Error() string {
	return fmt.Sprintf("branch not found: %s", e.Message)
}

func NewBranchNotExistsError(message string) *BranchNotExistsError {
	return &BranchNotExistsError{Message: message}
}
