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

	Fail(ctx context.Context, taskID uuid.UUID, err error) error
	Succeed(ctx context.Context, taskID uuid.UUID, code int, body Payload) error
}

type Task struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Status Status
	Error  *string
	Code   *int     // http status
	Body   *Payload // http body (JSON)

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
