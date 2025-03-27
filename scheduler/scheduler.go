package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"

	"github.com/google/uuid"
)

type Scheduler struct {
	taskRepo task.Repo
}

func New(ctx context.Context, taskRepo task.Repo) (*Scheduler, error) {
	s := &Scheduler{taskRepo: taskRepo}

	pendingTasks, err := s.taskRepo.ListPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("scheduler.New: %w", err)
	}

	for _, taskID := range pendingTasks {
		go s.Schedule(ctx, taskID)
	}

	return s, nil
}

func (s *Scheduler) Schedule(ctx context.Context, taskID uuid.UUID) error {
	t, err := s.taskRepo.Get(ctx, taskID)
	if err != nil {
		return s.taskRepo.Fail(ctx, taskID, err)
	}

	payload, _ := json.Marshal(t.Payload)
	res, err := http.Post(
		t.Target.String(),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return s.taskRepo.Fail(ctx, taskID, err)
	}
	defer res.Body.Close()

	var body task.Payload
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return s.taskRepo.Fail(ctx, taskID, err)
	}

	err = s.taskRepo.Succeed(ctx, taskID, res.StatusCode, body)
	if err != nil {
		return s.taskRepo.Fail(ctx, taskID, err)
	}

	return nil
}
