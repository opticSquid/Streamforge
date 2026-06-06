package config

import "time"

type Config struct {
	HlsRootDir       string
	RtmpBaseURi      string
	FFmpegBin        string
	FfmpegTimeoutSec time.Duration
}
