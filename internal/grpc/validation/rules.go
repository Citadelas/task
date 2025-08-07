package validation

import (
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/go-playground/validator/v10"
)

func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("task_priority", validateTaskPriority)

	v.RegisterValidation("task_status", validateTaskStatus)
}

func validateTaskPriority(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	_, exists := taskv1.TaskPriority_value[priority]
	return exists
}

func validateTaskStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()

	_, exists := taskv1.TaskStatus_value[status]
	if status == "TASK_STATUS_UNSPECIFIED" || status == "0" {
		return false
	}

	return exists
}
