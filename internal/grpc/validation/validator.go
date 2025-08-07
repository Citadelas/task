package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	validate *validator.Validate
)

func GetValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
		registerCustomValidators(validate)
	}
	return validate
}

func ValidateStruct(req interface{}) error {
	if err := GetValidator().Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(validationErrors)
		}
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	return nil
}

func formatValidationError(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param()))
		case "gt":
			messages = append(messages, fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param()))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s failed validation", err.Field()))
		}
	}
	return status.Errorf(codes.InvalidArgument, "validation errors: %v", messages)
}
