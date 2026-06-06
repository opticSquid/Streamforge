package transcoder

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func DefaultFFmpegConfig(binary, rtmpBase string, shutdownTimeout time.Duration) FFmpegConfig {
	return FFmpegConfig{
		Binary:          binary,
		RTMPBase:        rtmpBase,
		VideoScale:      "scale=-2:720",
		VideoBitrate:    "2500k",
		MaxRate:         "2500k",
		BufSize:         "5000k",
		GOPSize:         "60",
		AudioBitrate:    "128k",
		AudioRate:       "44100",
		HLSTime:         "2",
		HLSListSize:     "5",
		ShutdownTimeout: shutdownTimeout,
	}
}

func NewFFmpeg(cfg FFmpegConfig) *FFmpeg {
	return &FFmpeg{
		cfg: cfg,
	}
}

func (f *FFmpeg) buildArgs(name, outDir string) []string {
	return []string{
		"-i", f.cfg.RTMPBase + "/" + name,
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-tune", "zerolatency",
		"-vf", f.cfg.VideoScale,
		"-b:v", f.cfg.VideoBitrate,
		"-maxrate", f.cfg.MaxRate,
		"-bufsize", f.cfg.BufSize,
		"-g", f.cfg.GOPSize,
		"-keyint_min", f.cfg.GOPSize,
		"-sc_threshold", "0",
		"-c:a", "aac",
		"-b:a", f.cfg.AudioBitrate,
		"-ar", f.cfg.AudioRate,
		"-f", "hls",
		"-hls_time", f.cfg.HLSTime,
		"-hls_list_size", f.cfg.HLSListSize,
		"-hls_flags", "delete_segments+append_list",
		"-hls_segment_filename", filepath.Join(outDir, "seg_%03d.ts"),
		filepath.Join(outDir, "index.m3u8"),
	}
}

func (f *FFmpeg) Start(name, outDir string) (Process, error) {
	logPath := filepath.Join(outDir, "ffmpeg.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open ffmpeg log: %w", err)
	}

	args := f.buildArgs(name, outDir)
	cmd := exec.Command(f.cfg.Binary, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return nil, fmt.Errorf("start ffmpeg: %w", err)
	}

	proc := &ffmpegProcess{cmd: cmd, logFile: logFile}

	go func() {
		_ = proc.cmd.Wait()
		slog.Info("FFmpeg process exited", "name", name, "pid", cmd.Process.Pid)
	}()
	slog.Info("FFmpeg started", "name", name, "pid", cmd.Process.Pid)
	return proc, nil
}

func (p *ffmpegProcess) Stop(timeout time.Duration) {
	if p.cmd.Process == nil {
		return
	}
	pid := p.cmd.Process.Pid
	slog.Info("Sending SIGTERM to FFMPEG", "pid", pid)
	_ = p.cmd.Process.Signal(syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		_ = p.cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("FFMPEG exited cleanly", "pid", pid)
	case <-time.After(timeout):
		slog.Info("FFMPEG timed out, sending SIGKILL", "pid", pid)
		_ = p.cmd.Process.Signal(syscall.SIGKILL)
	}

	p.logFile.Close()
}
