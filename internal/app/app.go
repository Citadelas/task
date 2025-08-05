package app

import (
	grpcapp "github.com/Citadelas/task/internal/app/grpc"
	"github.com/Citadelas/task/internal/services/task"
	"github.com/Citadelas/task/internal/storage/postgresql"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	storage, err := postgresql.New(storagePath)
	if err != nil {
		panic(err)
	}
	taskService := task.New(log, storage, storage, storage, storage)
	grpcApp := grpcapp.New(log, taskService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
