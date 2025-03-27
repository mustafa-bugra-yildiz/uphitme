package userapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserApp struct {
	server *httptest.Server

	expectedHelloCount int
	gotHelloCount      int
}

func New(t *testing.T) *UserApp {
	t.Helper()

	ua := &UserApp{}

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/hello", ua.helloJob)

	server := httptest.NewServer(mux)
	ua.server = server

	return ua
}

func (ua *UserApp) Close(t *testing.T) {
	if ua.expectedHelloCount != ua.gotHelloCount {
		t.Fatalf("expected hello count was %d, got %d", ua.expectedHelloCount, ua.gotHelloCount)
	}
	ua.server.Close()
}

func (ua *UserApp) URL() string {
	return ua.server.URL
}

func (ua *UserApp) ExpectHello() {
	ua.expectedHelloCount++
}

func (ua *UserApp) helloJob(w http.ResponseWriter, r *http.Request) {
	ua.gotHelloCount++

	var payload struct {
		Message string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	got := payload.Message
	want := "hello world"
	if got != want {
		http.Error(
			w,
			fmt.Sprintf("/jobs/hello:\ngot:  %s\nwant: %s", got, want),
			http.StatusBadRequest,
		)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status": "okay",
	})
}
