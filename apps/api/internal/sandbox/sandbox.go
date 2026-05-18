package sandbox

import (
	"context"
	"errors"
	"time"
)

// Status represents the lifecycle state of a sandbox container.
type Status string

const (
	StatusCreating Status = "creating"
	StatusRunning  Status = "running"
	StatusStopped  Status = "stopped"
	StatusExited   Status = "exited"
	StatusUnknown  Status = "unknown"
)

// Sandbox is a managed Docker execution environment for an AI agent.
type Sandbox struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      Status    `json:"status"`
	Image       string    `json:"image"`
	ContainerID string    `json:"container_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
}

// CreateOptions configures a new sandbox instance.
type CreateOptions struct {
	Name  string            `json:"name,omitempty"`
	Image string            `json:"image,omitempty"`
	Env   map[string]string `json:"env,omitempty"`
	Cmd   []string          `json:"cmd,omitempty"`
}

// Manager defines Docker sandbox lifecycle operations.
type Manager interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, opts CreateOptions) (*Sandbox, error)
	Get(ctx context.Context, id string) (*Sandbox, error)
	List(ctx context.Context) ([]*Sandbox, error)
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
}

var (
	ErrNotFound      = errors.New("sandbox not found")
	ErrDockerUnavailable = errors.New("docker runtime unavailable")
)
