package validation

import (
	"github.com/go-playground/validator/v10"
)

func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("task_priority", validateTaskPriority)

	v.RegisterValidation("task_status", validateTaskStatus)
}

func validateTaskPriority(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	validPriorities := []string{"LOW", "MEDIUM", "HIGH"}
	for _, valid := range validPriorities {
		if priority == valid {
			return true
		}
	}
	return false
}

func validateTaskStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"TODO", "IN_PROGRESS", "DONE"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}
