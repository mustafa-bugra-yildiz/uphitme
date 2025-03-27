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
			status, title, target, payload
		FROM tasks
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
		var targetString string
		var payloadBytes []byte

		err := rows.Scan(
			&task.ID, &task.CreatedAt, &task.UpdatedAt,
			&task.Status, &task.Title, &targetString, &payloadBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("task.Repo.List: %w", err)
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

	_, err = p.db.ExecContext(
		ctx, query,
		id, title, target.String(), payloadBytes,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("task.Repo.Create: %w", err)
	}

	return id, nil
}
