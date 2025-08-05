package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/Citadelas/task/internal/domain/models"
	"github.com/Citadelas/task/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"
	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) CreateTask(ctx context.Context, title, description string,
	priority string) (*models.Task, error) {
	const op = "storage.postgresql.CreateTask"
	var task models.Task
	err := pgxscan.Get(ctx, s.db, &task, "INSERT INTO tasks(title, description, priority) "+
		"VALUES ($1, $2, $3) RETURNING id, title, description, priority, COALESCE(status, '') as status, created_at, due_date", title, description, priority)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) GetTask(ctx context.Context, id uint64) (*models.Task, error) {
	const op = "storage.postgresql.GetTask"
	var task models.Task
	err := pgxscan.Get(ctx, s.db, &task, "SELECT id, title, description, priority, COALESCE(status, '') as status, created_at, due_date FROM tasks WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id uint64, title, description, priority string) (*models.Task, error) {
	const op = "storage.postgresql.GetTask"
	var task models.Task
	query := `
        UPDATE tasks 
        SET 
            title = COALESCE(NULLIF($2, ''), title),
            description = COALESCE(NULLIF($3, ''), description),
            priority = COALESCE(NULLIF($4, ''), priority)
        WHERE id = $1 RETURNING *
    `
	err := pgxscan.Get(ctx, s.db, &task, query, id, title, description, priority)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, id uint64, status string) (*models.Task, error) {
	const op = "storage.postgresql.UpdateStatus"
	var task models.Task
	err := pgxscan.Get(ctx, s.db, &task, "UPDATE tasks SET status = $1 WHERE id = $2 RETURNING *", status, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) DeleteTask(ctx context.Context, id uint64) error {
	const op = "storage.postgresql.DeleteTask"
	commandTag, err := s.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrTaskNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
	}
	return nil
}
