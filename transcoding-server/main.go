package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/opticSquid/Streamforge/transcoding-server/config"
	"github.com/opticSquid/Streamforge/transcoding-server/hooks"
	"github.com/opticSquid/Streamforge/transcoding-server/stream"
	"github.com/opticSquid/Streamforge/transcoding-server/transcoder"
)

func main() {
	config := config.LoadEnv()

	if err := os.MkdirAll(config.HlsRootDir, 0o755); err != nil {
		log.Fatalf("cannot create HLS root %q: %v", config.HlsRootDir, err)
	}

	ffmpegCfg := transcoder.DefaultFFmpegConfig(config.FFmpegBin, config.RtmpBaseURi, config.FfmpegTimeoutSec)

	tc := transcoder.NewFFmpeg(ffmpegCfg)

	mgr := stream.NewManager(tc, config.HlsRootDir, config.FfmpegTimeoutSec)

	h := hooks.New(mgr)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	slog.Info("Stream controller starting", "port", config.Port)
	if err := http.ListenAndServe(config.Port, mux); err != nil {
		log.Fatalf("server error %v", err)
	}
}
