package utils

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// e164Regex matches strict E.164 phone format: + followed by 7-15 digits.
// Examples: +6281234567890 (valid), 08112233445 (invalid — no +)
var e164Regex = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

// RegisterCustomValidators registers custom validators for Gin's binding engine.
// Must be called once during application startup before handling requests.
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("e164_strict", validateE164Strict)
	}
}

// validateE164Strict enforces strict E.164 international phone format.
// Requires: starts with +, followed by 7-15 digits, no spaces or dashes.
func validateE164Strict(fl validator.FieldLevel) bool {
	return e164Regex.MatchString(fl.Field().String())
}
