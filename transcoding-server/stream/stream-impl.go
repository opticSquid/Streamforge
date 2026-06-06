package stream

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/opticSquid/Streamforge/transcoding-server/transcoder"
)

func NewManager(t transcoder.Transcoder, hlsRoot string, shutdownTimeout time.Duration) *Manager {
	return &Manager{
		streams:         make(map[string]*entry),
		transcoder:      t,
		hlsRoot:         hlsRoot,
		shutdownTimeout: shutdownTimeout,
	}
}

func (m *Manager) Start(name string) error {
	m.mu.Lock()
	if _, exists := m.streams[name]; exists {
		m.mu.Unlock()
		return fmt.Errorf("stream %q is already active", name)
	}
	m.mu.Unlock()

	outDir := filepath.Join(m.hlsRoot, name)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory for %q: %w", name, err)
	}

	proc, err := m.transcoder.Start(name, outDir)
	if err != nil {
		return fmt.Errorf("transcode start for %q: %w", name, err)
	}

	m.mu.Lock()
	m.streams[name] = &entry{proc: proc, outDir: outDir}
	m.mu.Unlock()

	slog.Info("stream started", "name", name, "outDir", outDir)
	return nil
}

func (m *Manager) Stop(name string) {
	m.mu.Lock()
	e, ok := m.streams[name]
	if ok {
		delete(m.streams, name)
	}
	m.mu.Unlock()

	if !ok {
		slog.Warn("Stop called for unknown stream", "name", name)
		return
	}

	slog.Info("Stopping stream", "name", name)
	e.proc.Stop(m.shutdownTimeout)

	if err := os.RemoveAll(e.outDir); err != nil {
		slog.Error("Failed to remove HLS directory", "name", name, "path", e.outDir, "err", err)
	} else {
		slog.Info("HLS directory removed", "name", name, "path", e.outDir)
	}
}

// Active returns a snapshot of currently active stream names.
func (m *Manager) Active() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	names := make([]string, 0, len(m.streams))
	for name := range m.streams {
		names = append(names, name)
	}
	return names
}
