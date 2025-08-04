package task

import (
	"context"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/Citadelas/task/internal/domain/models"
	"github.com/Citadelas/task/internal/grpc/converter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Task interface {
	CreateTask(ctx context.Context, title, description string,
		priority string) (*models.Task, error)

	GetTask(ctx context.Context, id uint64) (*models.Task, error)
	UpdateTask(ctx context.Context, id uint64, title, description string,
		priority string) (*models.Task, error)

	DeleteTask(ctx context.Context, id uint64) error
	UpdateStatus(ctx context.Context, id uint64, status string) (*models.Task, error)
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

//TODO: implement validator
//TODO: better enum validating
//TODO: less code duplicating

func (s *serverAPI) CreateTask(
	ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	if req.GetDescription() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "description is required")
	}
	if req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}
	if req.GetPriority() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "priority is required")
	}
	_, ok := taskv1.TaskPriority_value[req.GetPriority().String()]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect priority")
	}
	priority := req.GetPriority().String()
	task, err := s.task.CreateTask(ctx, req.GetTitle(), req.GetDescription(), priority)
	//TODO: add various errors handlers
	if err != nil {
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
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	task, err := s.task.GetTask(ctx, req.GetId())
	//TODO: add various errors handlers
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	res, err := s.adapter.ToProto(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &taskv1.GetTaskResponse{Task: res}, nil
}

func (s *serverAPI) UpdateTask(
	ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	task, err := s.task.UpdateTask(ctx, req.GetId(), req.GetTitle(), req.GetDescription(), req.GetPriority().String())
	if err != nil {
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
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	err := s.task.DeleteTask(ctx, req.GetId())
	//TODO: add various errors handlers
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) UpdateStatus(
	ctx context.Context, req *taskv1.UpdateStatusRequest) (*taskv1.UpdateStatusResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	_, ok := taskv1.TaskStatus_value[req.GetStatus().String()]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect status")
	}
	task, err := s.task.UpdateStatus(ctx, req.GetId(), req.GetStatus().String())
	//TODO: add various errors handlers
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	res, err := s.adapter.ToProto(task)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &taskv1.UpdateStatusResponse{Task: res}, nil
}
