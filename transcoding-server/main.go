package main

import (
	"log"
	"os"

	"github.com/opticSquid/Streamforge/transcoding-server/config"
)

func main() {
	config := config.LoadEnv()

	if err := os.MkdirAll(config.HlsRootDir, 0o755); err != nil {
		log.Fatalf("could not find HLS root directory, tried creating it but failed. error: %v", err)
	}

}
