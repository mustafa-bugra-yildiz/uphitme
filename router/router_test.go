package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mustafa-bugra-yildiz/uphitme/auth"
	"github.com/mustafa-bugra-yildiz/uphitme/mocks/userapp"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/scheduler"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestHealthWorks(t *testing.T) {
	ctrl := gomock.NewController(t)
	taskRepo := task.NewMockRepo(ctrl)
	userRepo := user.NewMockRepo(ctrl)
	auth := auth.New(userRepo)

	taskRepo.EXPECT().ListPending(gomock.Any()).Times(1)
	scheduler, _ := scheduler.New(context.Background(), taskRepo)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	res := httptest.NewRecorder()
	New(auth, taskRepo, userRepo, scheduler).ServeHTTP(res, req)

	got, _ := io.ReadAll(res.Body)
	got = bytes.TrimSpace(got)

	want, _ := json.Marshal(map[string]string{"status": "ok"})
	if bytes.Compare(got, want) != 0 {
		t.Errorf("\ngot:  %q\nwant: %q", string(got), string(want))
	}
}

func TestImmediateSchedulingWorks(t *testing.T) {
	ctrl := gomock.NewController(t)
	taskRepo := task.NewMockRepo(ctrl)
	userRepo := user.NewMockRepo(ctrl)
	auth := auth.New(userRepo)

	taskRepo.EXPECT().ListPending(gomock.Any()).Times(1)
	scheduler_, _ := scheduler.New(context.Background(), taskRepo)

	userApp := userapp.New(t)
	defer userApp.Close(t)

	id := uuid.New()
	target, _ := url.Parse(userApp.URL() + "/jobs/hello")
	payload := map[string]any{
		"message": "hello world",
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/schedule",
		makePayload(t, map[string]any{
			"title":   "Call me back with \"hello world\"",
			"target":  target.String(),
			"payload": payload,
		}),
	)

	taskRepo.EXPECT().
		Create(
			gomock.Any(),
			gomock.Eq("Call me back with \"hello world\""),
			gomock.Eq(*target),
			gomock.Eq(payload),
		).
		Return(id, nil).
		Times(1)
	userApp.ExpectHello()

	taskRepo.EXPECT().
		Get(gomock.Any(), gomock.Eq(id)).Times(1).
		Return(
			&task.Task{
				ID:      id,
				Target:  *target,
				Payload: payload,
			},
			nil,
		)
	taskRepo.EXPECT().
		Succeed(
			gomock.Any(),
			gomock.Eq(id),
			gomock.Eq(http.StatusOK),
			gomock.Eq(map[string]any{"status": "okay"}),
		).
		Times(1)

	res := httptest.NewRecorder()
	New(auth, taskRepo, userRepo, scheduler_).ServeHTTP(res, req)

	{
		got := res.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got:  %q\nwant: %q", got, want)
		}
	}

	{
		var got Response[string]
		json.NewDecoder(res.Body).Decode(&got)

		if !got.Succeeded {
			t.Errorf("response didn't succeed")
		}

		want := "task scheduled, we will call you back"
		if *got.Data != want {
			t.Errorf("response data wasn't what we wanted\ngot:  %q\nwant: %q", *got.Data, want)
		}
	}

	// wait for the hello job to be called
	time.Sleep(time.Second)
}

func makePayload(t *testing.T, v any) io.Reader {
	t.Helper()

	marshaled, err := json.Marshal(v)
	if err != nil {
		t.Errorf("payload failed: %q", err)
	}

	return bytes.NewBuffer(marshaled)
}
