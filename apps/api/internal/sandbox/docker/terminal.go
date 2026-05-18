package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type resizeMsg struct {
	Type string `json:"type"`
	Cols uint   `json:"cols"`
	Rows uint   `json:"rows"`
}

// AttachTerminal opens an interactive shell in the sandbox container over conn.
func (m *Manager) AttachTerminal(ctx context.Context, sandboxID string, conn sandbox.TerminalConn) error {
	c, err := m.findContainer(ctx, sandboxID)
	if err != nil {
		return err
	}

	execID, err := m.createShellExec(ctx, c.ID)
	if err != nil {
		return err
	}

	attach, err := m.client.ContainerExecAttach(ctx, execID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return fmt.Errorf("exec attach: %w", err)
	}
	defer attach.Close()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	done := make(chan struct{})

	// Docker stdout/stderr (TTY multiplexed) -> WebSocket client
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 32*1024)
		for {
			select {
			case <-done:
				return
			default:
			}
			n, readErr := attach.Conn.Read(buf)
			if n > 0 {
				if wErr := conn.WriteMessage(sandbox.TerminalMsgData, buf[:n]); wErr != nil {
					errCh <- wErr
					return
				}
			}
			if readErr != nil {
				if readErr != io.EOF {
					errCh <- readErr
				}
				return
			}
		}
	}()

	// WebSocket client -> Docker stdin
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
			}
			msgType, data, readErr := conn.ReadMessage()
			if readErr != nil {
				if readErr != io.EOF {
					errCh <- readErr
				}
				return
			}

			switch msgType {
			case sandbox.TerminalMsgResize:
				var rm resizeMsg
				if json.Unmarshal(data, &rm) == nil && rm.Type == "resize" && rm.Cols > 0 && rm.Rows > 0 {
					_ = m.client.ContainerExecResize(ctx, execID, container.ResizeOptions{
						Width:  rm.Cols,
						Height: rm.Rows,
					})
				}
			case sandbox.TerminalMsgData:
				if len(data) > 0 {
					if _, wErr := attach.Conn.Write(data); wErr != nil {
						errCh <- wErr
						return
					}
				}
			}
		}
	}()

	var runErr error
	select {
	case runErr = <-errCh:
	case <-ctx.Done():
		runErr = ctx.Err()
	}

	close(done)
	attach.Close()
	wg.Wait()
	return runErr
}

func (m *Manager) createShellExec(ctx context.Context, containerID string) (string, error) {
	for _, shell := range [][]string{{"/bin/bash"}, {"/bin/sh"}} {
		execResp, err := m.client.ContainerExecCreate(ctx, containerID, types.ExecConfig{
			User:         "ubuntu",
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			Cmd:          shell,
		})
		if err == nil {
			return execResp.ID, nil
		}
	}
	return "", fmt.Errorf("create shell exec")
}
