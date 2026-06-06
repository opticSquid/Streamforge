package stream

import (
	"sync"
	"time"

	"github.com/opticSquid/Streamforge/transcoding-server/transcoder"
)

type entry struct {
	proc   transcoder.Process
	outDir string
}

type Manager struct {
	mu              sync.Mutex
	streams         map[string]*entry
	transcoder      transcoder.Transcoder
	hlsRoot         string
	shutdownTimeout time.Duration
}
