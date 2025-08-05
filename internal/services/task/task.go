package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Citadelas/task/internal/domain/models"
	"github.com/Citadelas/task/internal/lib/logger/sl"
	"github.com/Citadelas/task/internal/storage"
	"log/slog"
)

type Task struct {
	logger  *slog.Logger
	getter  TaskGetter
	creator TaskCreator
	updater TaskUpdater
	deleter TaskDeleter
}

var (
	ErrWrongId = errors.New("wrong id")
)

type TaskCreator interface {
	CreateTask(ctx context.Context, title, description string,
		priority string) (*models.Task, error)
}

type TaskGetter interface {
	GetTask(ctx context.Context, id uint64) (*models.Task, error)
}
type TaskUpdater interface {
	UpdateTask(ctx context.Context, id uint64, title, description string,
		priority string) (*models.Task, error)
	UpdateStatus(ctx context.Context, id uint64, status string) (*models.Task, error)
}

type TaskDeleter interface {
	DeleteTask(ctx context.Context, id uint64) error
}

func New(
	log *slog.Logger,
	getter TaskGetter,
	creator TaskCreator,
	updater TaskUpdater,
	deleter TaskDeleter) *Task {

	return &Task{
		logger:  log,
		getter:  getter,
		creator: creator,
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
	res, err := t.creator.CreateTask(ctx, title, description, priority)
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
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn("task not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, ErrWrongId)
		}
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
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn("task not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, ErrWrongId)
		}
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
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn("task not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, ErrWrongId)
		}
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
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn("task not found", sl.Err(err))
			return fmt.Errorf("%s: %w", op, ErrWrongId)
		}
		log.Error("failed to delete task", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
