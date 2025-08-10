package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/Citadelas/task/internal/domain/models"
	"github.com/Citadelas/task/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

const returning = "RETURNING id, user_id, title, description, priority, COALESCE(status, '') as status, created_at, due_date"

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"
	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) CreateTask(ctx context.Context, uid uint64, title, description string,
	priority string) (*models.Task, error) {
	const op = "storage.postgresql.CreateTask"
	var task models.Task
	err := pgxscan.Get(ctx, s.db, &task, "INSERT INTO tasks(user_id, title, description, priority) "+
		"VALUES ($1, $2, $3, $4)"+returning, uid, title, description, priority)
	if err != nil {
		if lerr := checkTooLongField(op, err); lerr != nil {
			return nil, lerr
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) GetTask(ctx context.Context, id uint64, uid uint64) (*models.Task, error) {
	const op = "storage.postgresql.GetTask"
	var task models.Task
	err := pgxscan.Get(ctx, s.db, &task, ""+
		"SELECT id, user_id, title, description, priority, "+
		"COALESCE(status, '') as status, created_at, "+
		"due_date FROM tasks WHERE id = $1 AND user_id = $2", id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id uint64, uid uint64, title, description string,
	priority string) (*models.Task, error) {
	const op = "storage.postgresql.GetTask"
	var task models.Task
	query := `
        UPDATE tasks 
        SET 
            title = COALESCE(NULLIF($3, ''), title),
            description = COALESCE(NULLIF($4, ''), description),
            priority = COALESCE(NULLIF($5, ''), priority)
        WHERE id = $1 AND user_id = $2
    ` + returning
	err := pgxscan.Get(ctx, s.db, &task, query, id, uid, title, description, priority)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
		}
		if lerr := checkTooLongField(op, err); lerr != nil {
			return nil, lerr
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, id uint64, uid uint64, status string) (*models.Task, error) {
	const op = "storage.postgresql.UpdateStatus"
	var task models.Task
	err := pgxscan.Get(ctx, s.db, &task, "UPDATE tasks SET status = $1 WHERE id = $2 AND user_id = $3 RETURNING *", status, id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &task, nil
}

func (s *Storage) DeleteTask(ctx context.Context, id uint64, uid uint64) error {
	const op = "storage.postgresql.DeleteTask"
	commandTag, err := s.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
	}
	return nil
}

func checkTooLongField(op string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "22001" {
		return fmt.Errorf(
			"%s: %w", op, storage.ErrInputTooLong)
	}
	return nil
}
