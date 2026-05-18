package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/config"
	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

const labelSandboxID = "runtimewall.sandbox.id"

// Manager implements sandbox.Manager using the Docker Engine API.
type Manager struct {
	client   *client.Client
	cfg      config.Config
	events sandbox.EventStore
}

// NewManager connects to the Docker daemon.
func NewManager(cfg config.Config, events sandbox.EventStore) (*Manager, error) {
	opts := []client.Opt{client.FromEnv, client.WithAPIVersionNegotiation()}
	if cfg.DockerHost != "" {
		opts = append(opts, client.WithHost(cfg.DockerHost))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("docker client: %w", err)
	}

	return &Manager{client: cli, cfg: cfg, events: events}, nil
}

// Close releases the Docker client connection.
func (m *Manager) Close() error {
	return m.client.Close()
}

// Ping verifies the Docker daemon is reachable.
func (m *Manager) Ping(ctx context.Context) error {
	_, err := m.client.Ping(ctx)
	return err
}

// Create starts a new sandbox container.
func (m *Manager) Create(ctx context.Context, opts sandbox.CreateOptions) (*sandbox.Sandbox, error) {
	id := uuid.NewString()
	image := opts.Image
	if image == "" {
		image = m.cfg.SandboxImage
	}

	name := opts.Name
	if name == "" {
		name = fmt.Sprintf("runtimewall-sandbox-%s", id[:8])
	}

	env := make([]string, 0, len(opts.Env))
	for k, v := range opts.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	labels := map[string]string{
		m.cfg.SandboxLabel: "true",
		labelSandboxID:     id,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: false,
	}

	var networking *network.NetworkingConfig
	if m.cfg.SandboxNetwork != "" {
		networking = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				m.cfg.SandboxNetwork: {},
			},
		}
	}

	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Env:          env,
		Cmd:          opts.Cmd,
		Labels:       labels,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, hostConfig, networking, nil, name)
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}

	if err := m.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		_ = m.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return nil, fmt.Errorf("start container: %w", err)
	}

	now := time.Now().UTC()
	return &sandbox.Sandbox{
		ID:          id,
		Name:        name,
		Status:      sandbox.StatusRunning,
		Image:       image,
		ContainerID: resp.ID,
		CreatedAt:   now,
		StartedAt:   &now,
	}, nil
}

// Get returns a sandbox by RuntimeWall ID.
func (m *Manager) Get(ctx context.Context, id string) (*sandbox.Sandbox, error) {
	containers, err := m.listManaged(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, sandbox.ErrNotFound
}

// List returns all managed sandboxes.
func (m *Manager) List(ctx context.Context) ([]*sandbox.Sandbox, error) {
	return m.listManaged(ctx)
}

// Stop stops a running sandbox container.
func (m *Manager) Stop(ctx context.Context, id string) error {
	c, err := m.findContainer(ctx, id)
	if err != nil {
		return err
	}

	timeout := 10
	return m.client.ContainerStop(ctx, c.ID, container.StopOptions{Timeout: &timeout})
}

// Remove stops and deletes a sandbox container.
func (m *Manager) Remove(ctx context.Context, id string) error {
	c, err := m.findContainer(ctx, id)
	if err != nil {
		return err
	}

	timeout := 5
	_ = m.client.ContainerStop(ctx, c.ID, container.StopOptions{Timeout: &timeout})
	if err := m.client.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true}); err != nil {
		return err
	}
	if m.events != nil {
		m.events.Clear(id)
	}
	return nil
}

func (m *Manager) findContainer(ctx context.Context, sandboxID string) (types.Container, error) {
	list, err := m.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: m.managedFilters(),
	})
	if err != nil {
		return types.Container{}, fmt.Errorf("list containers: %w", err)
	}

	for _, c := range list {
		if c.Labels[labelSandboxID] == sandboxID {
			return c, nil
		}
	}
	return types.Container{}, sandbox.ErrNotFound
}

func (m *Manager) listManaged(ctx context.Context) ([]*sandbox.Sandbox, error) {
	list, err := m.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: m.managedFilters(),
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	out := make([]*sandbox.Sandbox, 0, len(list))
	for _, c := range list {
		sb, err := m.summaryToSandbox(c)
		if err != nil {
			continue
		}
		out = append(out, sb)
	}
	return out, nil
}

func (m *Manager) managedFilters() filters.Args {
	f := filters.NewArgs()
	f.Add("label", m.cfg.SandboxLabel+"=true")
	return f
}

func (m *Manager) summaryToSandbox(c types.Container) (*sandbox.Sandbox, error) {
	id := c.Labels[labelSandboxID]
	if id == "" {
		return nil, fmt.Errorf("missing sandbox id label")
	}

	name := c.ID
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}
	image := c.Image
	if len(c.ImageID) > 0 && image == "" {
		image = c.ImageID
	}

	var startedAt *time.Time
	if c.State == "running" && c.Status != "" {
		t := time.Now().UTC()
		startedAt = &t
	}

	return &sandbox.Sandbox{
		ID:          id,
		Name:        name,
		Status:      mapDockerState(c.State),
		Image:       image,
		ContainerID: c.ID,
		CreatedAt:   time.Unix(c.Created, 0).UTC(),
		StartedAt:   startedAt,
	}, nil
}

func mapDockerState(state string) sandbox.Status {
	switch strings.ToLower(state) {
	case "running":
		return sandbox.StatusRunning
	case "created", "restarting":
		return sandbox.StatusCreating
	case "exited", "dead":
		return sandbox.StatusExited
	case "paused", "removing":
		return sandbox.StatusStopped
	default:
		return sandbox.StatusUnknown
	}
}
