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
		return nil, ErrUnknownStatus
	}
	priorityVal, ok := taskv1.TaskPriority_value[domainTask.Priority]
	if !ok {
		return nil, ErrUnknownPriority
	}
	return &taskv1.Task{
		Id:          domainTask.Id,
		Title:       domainTask.Title,
		Description: domainTask.Description,
		Priority:    taskv1.TaskPriority(priorityVal),
		Status:      taskv1.TaskStatus(statusVal),
		CreatedAt:   timestamppb.New(domainTask.CreatedAt),
		DueDate:     timestamppb.New(domainTask.DueDate),
	}, nil
}

func (a *TaskAdapter) ToDomain(protoTask *taskv1.Task) (*models.Task, error) {
	return &models.Task{
		Id:          protoTask.Id,
		Title:       protoTask.Title,
		Description: protoTask.Description,
		Status:      protoTask.String(),
		Priority:    protoTask.String(),
		CreatedAt:   protoTask.CreatedAt.AsTime(),
		DueDate:     protoTask.DueDate.AsTime(),
	}, nil
}
