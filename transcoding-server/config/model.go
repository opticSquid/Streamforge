package config

import "time"

type Config struct {
	HlsRootDir       string
	RtmpBaseURi      string
	FFmpegBin        string
	Port             string
	FfmpegTimeoutSec time.Duration
}
