package log

import (
	"context"
	"log/slog"
	"slices"
	"sync"
)

func NewMulti(logs ...Handler) *MultiHandler {
	return &MultiHandler{logs: logs}
}

type MultiHandler struct {
	mu   sync.RWMutex
	logs []Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, h := range m.logs {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var last error
	for _, h := range m.logs {
		if err := h.Handle(ctx, r); err != nil {
			last = err
		}
	}
	return last
}

func (m *MultiHandler) AddHandler(h Handler) {
	if h == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !slices.Contains(m.logs, h) {
		m.logs = append(m.logs, h)
	}
}

func (m *MultiHandler) RemoveHandler(h Handler) {
	if h == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = slices.DeleteFunc(m.logs, func(h2 Handler) bool {
		return h == h2
	})
}
