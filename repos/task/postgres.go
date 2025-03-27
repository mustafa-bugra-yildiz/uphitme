package task

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type postgres struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return &postgres{db: db}
}

func (p postgres) List(ctx context.Context, page, pageSize int) ([]Task, error) {
	query := `
		SELECT
			id, created_at, updated_at,
			status, error, code, body,
			title, target, payload
		FROM tasks
		ORDER BY created_at DESC
		OFFSET $1
		LIMIT $2
	`

	rows, err := p.db.QueryContext(ctx, query, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, fmt.Errorf("task.Repo.List: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var bodyBytes *[]byte
		var targetString string
		var payloadBytes []byte

		err := rows.Scan(
			&task.ID, &task.CreatedAt, &task.UpdatedAt,
			&task.Status, &task.Error, &task.Code, &bodyBytes,
			&task.Title, &targetString, &payloadBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("task.Repo.List: %w", err)
		}

		if bodyBytes != nil {
			err = json.Unmarshal(*bodyBytes, &task.Body)
			if err != nil {
				return nil, fmt.Errorf("task.Repo.List: %w", err)
			}
		}

		target, err := url.Parse(targetString)
		if err != nil {
			return nil, fmt.Errorf("task.Repo.List: %w", err)
		}
		task.Target = *target

		err = json.Unmarshal(payloadBytes, &task.Payload)
		if err != nil {
			return nil, fmt.Errorf("task.Repo.List: %w", err)
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("task.Repo.List: %w", err)
	}

	return tasks, nil
}

func (p postgres) Create(
	ctx context.Context,
	title string,
	target url.URL,
	payload Payload,
) (uuid.UUID, error) {
	id := uuid.New()

	query := `
		INSERT INTO tasks (id, title, target, payload)
		VALUES ($1, $2, $3, $4)
	`

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("task.Repo.Create: %w", err)
	}

	result, err := p.db.ExecContext(
		ctx, query,
		id, title, target.String(), payloadBytes,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("task.Repo.Create: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, fmt.Errorf("task.Repo.Create: %w", err)
	}
	if affected == 0 {
		return uuid.Nil, fmt.Errorf("task.Repo.Create: no rows affected by insert")
	}

	return id, nil
}

func (p postgres) Fail(
	ctx context.Context,
	taskID uuid.UUID,
	err error,
) error {
	query := `
		UPDATE tasks
		SET
			status = 'failed',
			error  = $2
		WHERE id = $1
	`

	result, err := p.db.ExecContext(ctx, query, taskID, err.Error())
	if err != nil {
		return fmt.Errorf("task.Repo.Fail: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("task.Repo.Fail: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("task.Repo.Fail: no rows affected by update")
	}

	return nil
}

func (p postgres) Succeed(
	ctx context.Context,
	taskID uuid.UUID,
	code int,
	body Payload,
) error {
	query := `
		UPDATE tasks
		SET
			status = 'succeeded',
			code   = $2
			body   = $3
		WHERE id = $1
	`

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("task.Repo.Succeed: %w", err)
	}

	result, err := p.db.ExecContext(ctx, query, taskID, code, bodyBytes)
	if err != nil {
		return fmt.Errorf("task.Repo.Succeed: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("task.Repo.Succeed: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("task.Repo.Succeed: no rows affected by update")
	}

	return nil
}
