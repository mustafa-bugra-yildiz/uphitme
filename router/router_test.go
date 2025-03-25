package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthWorks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	got, _ := io.ReadAll(res.Body)
	got = bytes.TrimSpace(got)

	want, _ := json.Marshal(map[string]string{"status": "ok"})
	if bytes.Compare(got, want) != 0 {
		t.Errorf("\ngot:  %q\nwant: %q", string(got), string(want))
	}
}

func TestImmediateSchedulingWorks(t *testing.T) {
	handler := http.NewServeMux()

	f, isCalled := helloJob(t)
	handler.HandleFunc("/jobs/hello", f)

	server := httptest.NewServer(handler)
	defer server.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/schedule",
		payload(t, map[string]any{
			"title":  "Call me back with \"hello world\"",
			"target": server.URL + "/jobs/hello",
			"payload": map[string]any{
				"message": "hello world",
			},
		}),
	)

	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

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

	time.Sleep(time.Second)
	if !*isCalled {
		t.Errorf("/jobs/hello was never called")
	}
}

func helloJob(t *testing.T) (func(w http.ResponseWriter, r *http.Request), *bool) {
	isCalledValue := false
	isCalled := &isCalledValue

	handler := func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		var payload struct {
			Message string `json:"message"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Errorf("/jobs/hello: %q", err)
		}

		got := payload.Message
		want := "hello world"

		if got != want {
			t.Errorf("/jobs/hello:\ngot:  %s\nwant: %s", got, want)
		}

		*isCalled = true
		w.Write([]byte("Gotta return 200"))
	}

	return handler, isCalled
}

func payload(t *testing.T, v any) io.Reader {
	t.Helper()

	marshaled, err := json.Marshal(v)
	if err != nil {
		t.Errorf("payload failed: %q", err)
	}

	return bytes.NewBuffer(marshaled)
}
