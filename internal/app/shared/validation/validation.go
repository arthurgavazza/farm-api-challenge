package shared

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

func parseValidationError(errs error) []ErrorResponse {
	validationErrors := []ErrorResponse{}
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}

var validate = validator.New()

func ValidateStruct(data interface{}) []ErrorResponse {
	errs := validate.Struct(data)
	return parseValidationError(errs)
}
