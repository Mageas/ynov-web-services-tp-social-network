package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validator provides validation utilities
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if there are no validation errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// GetErrors returns the map of validation errors
func (v *Validator) GetErrors() map[string]string {
	return v.Errors
}

// AddError adds an error message for a given field
func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// Check adds an error if the condition is not met
func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// Required checks if a value is not empty
func (v *Validator) Required(value string, field string) {
	v.Check(strings.TrimSpace(value) != "", field, "this field is required")
}

// Email checks if a value is a valid email
func (v *Validator) Email(value string, field string) {
	v.Check(emailRegex.MatchString(value), field, "must be a valid email address")
}

// MinLength checks if a value has at least n characters
func (v *Validator) MinLength(value string, n int, field string) {
	v.Check(utf8.RuneCountInString(value) >= n, field, fmt.Sprintf("must be at least %d characters", n))
}

// MaxLength checks if a value has at most n characters
func (v *Validator) MaxLength(value string, n int, field string) {
	v.Check(utf8.RuneCountInString(value) <= n, field, fmt.Sprintf("must be at most %d characters", n))
}

// Range checks if a value is between min and max characters
func (v *Validator) Range(value string, min, max int, field string) {
	length := utf8.RuneCountInString(value)
	v.Check(length >= min && length <= max, field, fmt.Sprintf("must be between %d and %d characters", min, max))
}
