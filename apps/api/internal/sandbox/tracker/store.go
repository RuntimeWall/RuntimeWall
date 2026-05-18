package tracker

import (
	"log/slog"
	"sync"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/security"
	"github.com/google/uuid"
)

// Store is the in-memory runtime event and policy store.
type Store struct {
	mu       sync.RWMutex
	events   map[string][]sandbox.RuntimeEvent
	subs     map[string]map[chan sandbox.RuntimeEvent]struct{}
	policies map[string]sandbox.SecurityPolicy
}

// NewStore creates an event store.
func NewStore() *Store {
	return &Store{
		events:   make(map[string][]sandbox.RuntimeEvent),
		subs:     make(map[string]map[chan sandbox.RuntimeEvent]struct{}),
		policies: make(map[string]sandbox.SecurityPolicy),
	}
}

// RecordCommand classifies, enforces policy, and stores a runtime event.
func (s *Store) RecordCommand(sandboxID, command string, source sandbox.CommandSource) sandbox.RuntimeEvent {
	cls := security.Classify(command)
	policy := s.GetPolicy(sandboxID)
	blocked, blockReason := security.EvaluatePolicy(policy, cls, command)

	reason := cls.Reason
	if blocked {
		reason = blockReason
	}

	ev := sandbox.RuntimeEvent{
		ID:        uuid.NewString(),
		SandboxID: sandboxID,
		Event:     cls.EventType,
		Command:   command,
		Threat:    cls.Threat,
		Blocked:   blocked,
		Reason:    reason,
		Source:    source,
		Timestamp: time.Now().UTC(),
	}

	if blocked {
		ev.Event = sandbox.EventPolicyViolation
		slog.Warn("runtime policy blocked command",
			"sandbox_id", sandboxID,
			"command", command,
			"threat", cls.Threat,
			"reason", blockReason,
		)
	}

	s.publish(ev)
	return ev
}

func (s *Store) publish(ev sandbox.RuntimeEvent) {
	s.mu.Lock()
	s.events[ev.SandboxID] = append(s.events[ev.SandboxID], ev)
	chans := make([]chan sandbox.RuntimeEvent, 0, len(s.subs[ev.SandboxID]))
	for ch := range s.subs[ev.SandboxID] {
		chans = append(chans, ch)
	}
	s.mu.Unlock()

	for _, ch := range chans {
		select {
		case ch <- ev:
		default:
		}
	}
}

// List returns events for a sandbox.
func (s *Store) List(sandboxID string) []sandbox.RuntimeEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.events[sandboxID]
	if len(items) == 0 {
		return []sandbox.RuntimeEvent{}
	}
	out := make([]sandbox.RuntimeEvent, len(items))
	copy(out, items)
	return out
}

// Subscribe listens for new events.
func (s *Store) Subscribe(sandboxID string) (<-chan sandbox.RuntimeEvent, func()) {
	ch := make(chan sandbox.RuntimeEvent, 32)
	s.mu.Lock()
	if s.subs[sandboxID] == nil {
		s.subs[sandboxID] = make(map[chan sandbox.RuntimeEvent]struct{})
	}
	s.subs[sandboxID][ch] = struct{}{}
	s.mu.Unlock()

	return ch, func() {
		s.mu.Lock()
		delete(s.subs[sandboxID], ch)
		close(ch)
		s.mu.Unlock()
	}
}

// Clear removes sandbox events.
func (s *Store) Clear(sandboxID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, sandboxID)
	delete(s.subs, sandboxID)
	delete(s.policies, sandboxID)
}

// GetPolicy returns policy for a sandbox (default if unset).
func (s *Store) GetPolicy(sandboxID string) sandbox.SecurityPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if p, ok := s.policies[sandboxID]; ok {
		return p
	}
	return sandbox.DefaultPolicy()
}

// SetPolicy sets per-sandbox policy.
func (s *Store) SetPolicy(sandboxID string, policy sandbox.SecurityPolicy) {
	s.mu.Lock()
	s.policies[sandboxID] = policy
	s.mu.Unlock()
}

// Record implements legacy CommandTracker for compatibility.
func (s *Store) Record(sandboxID, command string, source sandbox.CommandSource) {
	s.RecordCommand(sandboxID, command, source)
}
