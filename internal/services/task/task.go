package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Citadelas/task/internal/domain/models"
	"github.com/Citadelas/task/internal/lib/logger/sl"
	"log/slog"
)

type Task struct {
	logger  *slog.Logger
	getter  Getter
	creater Creater
	updater Updater
	deleter Deleter
}

var (
	ErrWrongId = errors.New("wrong id")
)

type Creater interface {
	CreateTask(ctx context.Context, title, description string,
		priority string) (*models.Task, error)
}

type Getter interface {
	GetTask(ctx context.Context, id uint64) (*models.Task, error)
}

type Updater interface {
	UpdateTask(ctx context.Context, id uint64, title, description string,
		priority string) (*models.Task, error)
	UpdateStatus(ctx context.Context, id uint64, status string) (*models.Task, error)
}

type Deleter interface {
	DeleteTask(ctx context.Context, id uint64) error
}

func New(
	log *slog.Logger,
	getter Getter,
	creater Creater,
	updater Updater,
	deleter Deleter) *Task {

	return &Task{
		logger:  log,
		getter:  getter,
		creater: creater,
		updater: updater,
		deleter: deleter,
	}
}

//TODO: better error handling
//TODO: better interface naming and location

func (t *Task) CreateTask(ctx context.Context, title, description string,
	priority string) (*models.Task, error) {
	const op = "task.CreateTask"
	log := t.logger.With(
		slog.String("op", op),
	)
	res, err := t.creater.CreateTask(ctx, title, description, priority)
	if err != nil {
		log.Error("Failed to create task", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (t *Task) GetTask(ctx context.Context, id uint64) (*models.Task, error) {
	const op = "task.GetTask"
	log := t.logger.With(
		slog.String("op", op),
	)
	res, err := t.getter.GetTask(ctx, id)
	if err != nil {
		log.Error("failed to get task", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (t *Task) UpdateTask(ctx context.Context, id uint64, title, description string,
	priority string) (*models.Task, error) {
	const op = "task.UpdateTask"
	log := t.logger.With(
		slog.String("op", op),
	)
	res, err := t.updater.UpdateTask(ctx, id, title, description, priority)
	if err != nil {
		log.Error("failed to update task", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (t *Task) UpdateStatus(ctx context.Context, id uint64, status string) (*models.Task, error) {
	const op = "task.UpdateTask"
	log := t.logger.With(
		slog.String("op", op),
	)
	res, err := t.updater.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error("failed to update status", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (t *Task) DeleteTask(ctx context.Context, id uint64) error {
	const op = "task.UpdateTask"
	log := t.logger.With(
		slog.String("op", op),
	)
	err := t.deleter.DeleteTask(ctx, id)
	if err != nil {
		log.Error("failed to delete task", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
