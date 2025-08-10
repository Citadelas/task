package converter

import (
	"errors"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/Citadelas/task/internal/domain/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUnknownStatus   = errors.New("unknown task status")
	ErrUnknownPriority = errors.New("unknown task priority")
)

type TaskAdapter struct{}

func NewTaskAdapter() *TaskAdapter {
	return &TaskAdapter{}
}

func (a *TaskAdapter) ToProto(domainTask *models.Task) (*taskv1.Task, error) {
	statusVal, ok := taskv1.TaskStatus_value[domainTask.Status]
	if !ok {
		if domainTask.Status != "" {
			return nil, ErrUnknownStatus
		}
		statusVal = taskv1.TaskStatus_value["TASK_STATUS_UNSPECIFIED"]
	}
	status := taskv1.TaskStatus(statusVal)

	priorityVal, ok := taskv1.TaskPriority_value[domainTask.Priority]
	if !ok {
		return nil, ErrUnknownPriority
	}
	return &taskv1.Task{
		Id:          domainTask.Id,
		UserId:      domainTask.UserId,
		Title:       domainTask.Title,
		Description: domainTask.Description,
		Priority:    taskv1.TaskPriority(priorityVal),
		Status:      status,
		CreatedAt:   timestamppb.New(domainTask.CreatedAt),
		DueDate:     timestamppb.New(domainTask.DueDate),
	}, nil
}

func (a *TaskAdapter) ToDomain(protoTask *taskv1.Task) (*models.Task, error) {
	return &models.Task{
		Id:          protoTask.Id,
		UserId:      protoTask.UserId,
		Title:       protoTask.Title,
		Description: protoTask.Description,
		Status:      protoTask.String(),
		Priority:    protoTask.String(),
		CreatedAt:   protoTask.CreatedAt.AsTime(),
		DueDate:     protoTask.DueDate.AsTime(),
	}, nil
}
