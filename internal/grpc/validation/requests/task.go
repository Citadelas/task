package requests

type CreateTaskRequest struct {
	UID         uint64 `validate:"required,gt=0"`
	Title       string `validate:"required,min=1,max=200"`
	Description string `validate:"required,min=1,max=1000"`
	Priority    string `validate:"required,oneof=LOW MEDIUM HIGH"`
}

type GetTaskRequest struct {
	ID  uint64 `validate:"required,gt=0"`
	UID uint64 `validate:"required,gt=0"`
}

type UpdateTaskRequest struct {
	ID          uint64 `validate:"required,gt=0"`
	UID         uint64 `validate:"required,gt=0"`
	Title       string `validate:"omitempty,min=1,max=200"`
	Description string `validate:"omitempty,min=1,max=1000"`
	Priority    string `validate:"omitempty,oneof=LOW MEDIUM HIGH"`
}

type DeleteTaskRequest struct {
	ID  uint64 `validate:"required,gt=0"`
	UID uint64 `validate:"required,gt=0"`
}

type UpdateStatusRequest struct {
	ID     uint64 `validate:"required,gt=0"`
	UID    uint64 `validate:"required,gt=0"`
	Status string `validate:"required,task_status"`
}
