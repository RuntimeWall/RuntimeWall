package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/google/uuid"
)

// CreateUbuntu launches an isolated Ubuntu sandbox container.
func (m *Manager) CreateUbuntu(ctx context.Context) (*sandbox.LaunchResult, error) {
	imageRef := m.cfg.UbuntuSandboxImage
	if err := m.ensureImage(ctx, imageRef); err != nil {
		return nil, fmt.Errorf("ensure image %s: %w", imageRef, err)
	}

	sandboxID := uuid.NewString()
	name := fmt.Sprintf("runtimewall-ubuntu-%s", sandboxID[:8])

	labels := map[string]string{
		m.cfg.SandboxLabel: "true",
		labelSandboxID:     sandboxID,
		"runtimewall.os":   "ubuntu",
	}

	hostConfig := &container.HostConfig{
		AutoRemove:     false,
		ReadonlyRootfs: true,
		CapDrop:        []string{"ALL"},
		SecurityOpt:    []string{"no-new-privileges:true"},
		NetworkMode:    container.NetworkMode("none"),
		Tmpfs: map[string]string{
			"/tmp": "rw,noexec,nosuid,size=64m",
			"/run": "rw,noexec,nosuid,size=32m",
		},
		Resources: container.Resources{
			Memory:   512 * 1024 * 1024,
			NanoCPUs: 500_000_000,
		},
	}

	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image:  imageRef,
		Cmd:    []string{"sleep", "infinity"},
		Labels: labels,
		User:   "1000:1000",
	}, hostConfig, nil, nil, name)
	if err != nil {
		return nil, fmt.Errorf("create ubuntu sandbox: %w", err)
	}

	if err := m.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		_ = m.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return nil, fmt.Errorf("start ubuntu sandbox: %w", err)
	}

	return &sandbox.LaunchResult{
		ContainerID: resp.ID,
		Image:       imageRef,
		Status:      sandbox.StatusRunning,
	}, nil
}

func (m *Manager) ensureImage(ctx context.Context, ref string) error {
	_, _, err := m.client.ImageInspectWithRaw(ctx, ref)
	if err == nil {
		return nil
	}

	pullCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	reader, err := m.client.ImagePull(pullCtx, ref, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	_, _ = io.Copy(io.Discard, reader)
	return nil
}
