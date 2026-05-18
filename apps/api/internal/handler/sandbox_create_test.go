package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

func TestCreateUbuntu(t *testing.T) {
	h := NewSandboxes(stubManager{})
	req := httptest.NewRequest(http.MethodPost, "/sandbox/create", nil)
	rec := httptest.NewRecorder()

	h.CreateUbuntu(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var body sandbox.LaunchResult
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.ContainerID != "abc123" {
		t.Fatalf("container_id = %q", body.ContainerID)
	}
	if body.Image != "ubuntu:22.04" {
		t.Fatalf("image = %q", body.Image)
	}
	if body.Status != sandbox.StatusRunning {
		t.Fatalf("status = %q", body.Status)
	}
}

func TestCreateUbuntu_NoDocker(t *testing.T) {
	h := NewSandboxes(nil)
	req := httptest.NewRequest(http.MethodPost, "/sandbox/create", nil)
	rec := httptest.NewRecorder()

	h.CreateUbuntu(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
