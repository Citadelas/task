package TaskService

import (
	"context"
	"errors"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/Citadelas/task/internal/domain/models"
	"github.com/Citadelas/task/internal/grpc/converter"
	"github.com/Citadelas/task/internal/grpc/validation"
	"github.com/Citadelas/task/internal/grpc/validation/requests"
	taskservice "github.com/Citadelas/task/internal/services/task"
	"github.com/Citadelas/task/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Task interface {
	CreateTask(ctx context.Context, uid uint64,
		title, description, priority string) (*models.Task, error)

	GetTask(ctx context.Context, id, uid uint64) (*models.Task, error)
	UpdateTask(ctx context.Context, id, uid uint64, title, description string,
		priority string) (*models.Task, error)

	DeleteTask(ctx context.Context, id, uid uint64) error
	UpdateStatus(ctx context.Context, id, uid uint64, status string) (*models.Task, error)
}

type serverAPI struct {
	task    Task
	adapter *converter.TaskAdapter
	taskv1.UnimplementedTaskServiceServer
}

func Register(gRPC *grpc.Server, task Task) {
	taskv1.RegisterTaskServiceServer(gRPC, &serverAPI{
		task:    task,
		adapter: converter.NewTaskAdapter(),
	})
}

func (s *serverAPI) CreateTask(
	ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	validationReq := requests.CreateTaskRequest{
		UID:         req.GetUserId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Priority:    req.GetPriority().String(),
	}
	if err := validation.ValidateStruct(validationReq); err != nil {
		return nil, err
	}
	priority := req.GetPriority().String()
	task, err := s.task.CreateTask(ctx, req.GetUserId(), req.GetTitle(), req.GetDescription(), priority)
	if err != nil {
		if errors.Is(err, storage.ErrInputTooLong) {
			return nil, status.Error(codes.InvalidArgument, storage.ErrInputTooLong.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	res, err := s.adapter.ToProto(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &taskv1.CreateTaskResponse{Task: res}, err
}

func (s *serverAPI) GetTask(
	ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error) {

	validationReq := requests.GetTaskRequest{ID: req.GetId(), UID: req.GetUserId()}
	if err := validation.ValidateStruct(validationReq); err != nil {
		return nil, err
	}

	task, err := s.task.GetTask(ctx, req.GetId(), req.GetUserId())
	if err != nil {
		if errors.Is(err, taskservice.ErrWrongId) {
			return nil, status.Error(codes.InvalidArgument, "task not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	res, err := s.adapter.ToProto(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &taskv1.GetTaskResponse{Task: res}, nil
}

func (s *serverAPI) UpdateTask(
	ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error) {
	validationReq := requests.UpdateTaskRequest{
		ID:          req.GetId(),
		UID:         req.GetUserId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Priority:    req.GetPriority().String(),
	}
	if err := validation.ValidateStruct(validationReq); err != nil {
		return nil, err
	}

	task, err := s.task.UpdateTask(ctx, req.GetId(), req.GetUserId(), req.GetTitle(), req.GetDescription(), req.GetPriority().String())
	if err != nil {
		if errors.Is(err, taskservice.ErrWrongId) {
			return nil, status.Error(codes.InvalidArgument, "task not found")
		}
		if errors.Is(err, storage.ErrInputTooLong) {
			return nil, status.Error(codes.InvalidArgument, storage.ErrInputTooLong.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	res, err := s.adapter.ToProto(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &taskv1.UpdateTaskResponse{Task: res}, nil
}

func (s *serverAPI) DeleteTask(
	ctx context.Context, req *taskv1.DeleteTaskRequest) (*emptypb.Empty, error) {

	validationReq := requests.DeleteTaskRequest{
		ID:  req.GetId(),
		UID: req.GetUserId(),
	}
	if err := validation.ValidateStruct(validationReq); err != nil {
		return nil, err
	}

	err := s.task.DeleteTask(ctx, req.GetId(), req.GetUserId())
	if err != nil {
		if errors.Is(err, taskservice.ErrWrongId) {
			return nil, status.Error(codes.InvalidArgument, "task not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) UpdateStatus(
	ctx context.Context, req *taskv1.UpdateStatusRequest) (*taskv1.UpdateStatusResponse, error) {
	validationReq := requests.UpdateStatusRequest{
		ID:     req.GetId(),
		Status: req.GetStatus().String(),
		UID:    req.GetUserId(),
	}
	if err := validation.ValidateStruct(validationReq); err != nil {
		return nil, err
	}

	task, err := s.task.UpdateStatus(ctx, req.GetId(), req.GetUserId(), req.GetStatus().String())
	if err != nil {
		if errors.Is(err, taskservice.ErrWrongId) {
			return nil, status.Error(codes.InvalidArgument, "task not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	res, err := s.adapter.ToProto(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &taskv1.UpdateStatusResponse{Task: res}, nil
}
