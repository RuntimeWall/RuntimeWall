package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

type stubManager struct{}

func (stubManager) Ping(context.Context) error { return nil }
func (stubManager) CreateUbuntu(context.Context) (*sandbox.LaunchResult, error) {
	return &sandbox.LaunchResult{
		ContainerID: "abc123",
		Image:       "ubuntu:22.04",
		Status:      sandbox.StatusRunning,
	}, nil
}
func (stubManager) Create(context.Context, sandbox.CreateOptions) (*sandbox.Sandbox, error) {
	return nil, nil
}
func (stubManager) Get(context.Context, string) (*sandbox.Sandbox, error) { return nil, nil }
func (stubManager) List(context.Context) ([]*sandbox.Sandbox, error)      { return nil, nil }
func (stubManager) Stop(context.Context, string) error                    { return nil }
func (stubManager) Remove(context.Context, string) error                  { return nil }

func TestHealth_OK(t *testing.T) {
	h := NewHealth(stubManager{}, "0.1.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %v, want ok", body["status"])
	}
}

func TestHealth_NoDocker(t *testing.T) {
	h := NewHealth(nil, "0.1.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
