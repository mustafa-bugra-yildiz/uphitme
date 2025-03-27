package scheduler

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/mustafa-bugra-yildiz/uphitme/mocks/userapp"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestImmediateSchedulingWorks(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	taskRepo := task.NewMockRepo(ctrl)

	taskRepo.EXPECT().ListPending(gomock.Eq(ctx)).Times(1)
	scheduler, _ := New(ctx, taskRepo)

	userApp := userapp.New(t)
	defer userApp.Close(t)

	id := uuid.New()
	target, _ := url.Parse(userApp.URL() + "/jobs/hello")
	payload := map[string]any{
		"message": "hello world",
	}

	userApp.ExpectHello()
	taskRepo.EXPECT().
		Get(gomock.Eq(ctx), gomock.Eq(id)).
		Return(
			&task.Task{
				ID:      id,
				Target:  *target,
				Payload: payload,
			},
			nil,
		).
		Times(1)
	taskRepo.EXPECT().
		Succeed(
			gomock.Eq(ctx),
			gomock.Eq(id),
			gomock.Eq(http.StatusOK),
			gomock.Eq(map[string]any{
				"status": "okay",
			}),
		).
		Times(1)

	err := scheduler.Schedule(ctx, id)
	if err != nil {
		t.Fatalf("error scheduling task: %q", err)
	}
}
