package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Invalid port, needs an integer value")
	}
	ffmpegTimeoutSec, err := strconv.Atoi(os.Getenv("FFMPEG_TIMEOUT_SEC"))
	if err != nil {
		log.Fatal("Invalid FFMPEG_TIMEOUT_SEC, needs an integer value")
	}

	return Config{
		HlsRootDir:       os.Getenv("HLS_ROOT_DIR"),
		RtmpBaseURi:      os.Getenv("RTMP_BASE_URI"),
		Port:             port,
		FFmpegBin:        os.Getenv("FFMPEG_BINARY"),
		FfmpegTimeoutSec: time.Duration(ffmpegTimeoutSec) * time.Second,
	}
}

func Default() Config {
	return Config{
		HlsRootDir:       "./hls",
		RtmpBaseURi:      "rtmp://localhost:1935/live",
		FFmpegBin:        "ffmpeg",
		Port:             8081,
		FfmpegTimeoutSec: 5 * time.Second,
	}
}
