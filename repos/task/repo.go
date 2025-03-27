package task

import (
	"context"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	List(ctx context.Context, page, pageSize int) ([]Task, error)
	Create(
		ctx context.Context,
		title string,
		target url.URL,
		payload Payload,
	) (uuid.UUID, error)
	SetStatus(ctx context.Context, taskID uuid.UUID, status Status) error
}

type Task struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Status  Status
	Title   string
	Target  url.URL
	Payload Payload
}

type Status string

const (
	Pending   Status = "pending"
	Failed    Status = "failed"
	Succeeded Status = "succeeded"
)

type Payload map[string]any
