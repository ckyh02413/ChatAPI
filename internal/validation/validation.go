package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		var msgs []string
		for _, e := range err.(validator.ValidationErrors) {
			msgs = append(msgs, formatError(e))
		}
		return fmt.Errorf("%s", strings.Join(msgs, ", "))
	}
	return nil
}

func formatError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
