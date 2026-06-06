package transcoder

import (
	"os"
	"os/exec"
	"time"
)

type Process interface {
	Stop(timeout time.Duration)
}

type Transcoder interface {
	Start(name, outDir string) (Process, error)
}

type ffmpegProcess struct {
	cmd     *exec.Cmd
	logFile *os.File
}

type FFmpegConfig struct {
	Binary          string
	RTMPBase        string
	VideoScale      string
	VideoBitrate    string
	MaxRate         string
	BufSize         string
	GOPSize         string
	AudioBitrate    string
	AudioRate       string
	HLSTime         string
	HLSListSize     string
	ShutdownTimeout time.Duration
}

type FFmpeg struct {
	cfg FFmpegConfig
}
