package docker

import (
	"context"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types/container"
)

// Sweeper periodically reaps stopped sandbox containers older than ttl. It is
// safe to call Start on a nil Manager (it becomes a no-op).
type Sweeper struct {
	mgr      *Manager
	ttl      time.Duration
	interval time.Duration
	cancel   context.CancelFunc
}

// NewSweeper builds a sweeper. If ttl <= 0 the sweeper is disabled and
// Start/Stop are no-ops.
func NewSweeper(mgr *Manager, ttl, interval time.Duration) *Sweeper {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Sweeper{mgr: mgr, ttl: ttl, interval: interval}
}

// Start begins the background sweep loop. Returns immediately.
func (s *Sweeper) Start(parent context.Context) {
	if s == nil || s.mgr == nil || s.ttl <= 0 {
		return
	}
	ctx, cancel := context.WithCancel(parent)
	s.cancel = cancel

	slog.Info("sandbox sweeper started", "ttl", s.ttl.String(), "interval", s.interval.String())

	go func() {
		t := time.NewTicker(s.interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				sweepCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				removed, err := s.mgr.SweepStopped(sweepCtx, s.ttl)
				cancel()
				if err != nil {
					slog.Warn("sandbox sweep failed", "error", err)
					continue
				}
				if len(removed) > 0 {
					slog.Info("sandbox sweep reaped containers", "count", len(removed), "ids", removed)
				}
			}
		}
	}()
}

// Stop halts the background sweep loop.
func (s *Sweeper) Stop() {
	if s == nil || s.cancel == nil {
		return
	}
	s.cancel()
}

// SweepStopped removes managed containers whose state is exited/dead/created
// and whose FinishedAt (or Created, if never started) is older than now-ttl.
// Returns the RuntimeWall sandbox IDs that were removed.
func (m *Manager) SweepStopped(ctx context.Context, ttl time.Duration) ([]string, error) {
	if ttl <= 0 {
		return nil, nil
	}

	list, err := m.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: m.managedFilters(),
	})
	if err != nil {
		return nil, err
	}

	removed := make([]string, 0)
	cutoff := time.Now().Add(-ttl)

	for _, c := range list {
		switch c.State {
		case "exited", "dead", "created":
			// reapable
		default:
			continue
		}

		sandboxID := c.Labels[labelSandboxID]
		if sandboxID == "" {
			continue
		}

		stoppedAt := time.Unix(c.Created, 0)
		info, err := m.client.ContainerInspect(ctx, c.ID)
		if err == nil && info.State != nil && info.State.FinishedAt != "" {
			if t, perr := time.Parse(time.RFC3339Nano, info.State.FinishedAt); perr == nil && !t.IsZero() {
				stoppedAt = t
			}
		}

		if stoppedAt.After(cutoff) {
			continue
		}

		if err := m.client.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true}); err != nil {
			slog.Warn("sweeper failed to remove container", "container", c.ID, "error", err)
			continue
		}
		if m.events != nil {
			m.events.Clear(sandboxID)
		}
		removed = append(removed, sandboxID)
		slog.Info("sweeper removed stopped sandbox",
			"sandbox_id", sandboxID,
			"container", c.ID,
			"stopped_at", stoppedAt.Format(time.RFC3339),
		)
	}
	return removed, nil
}
