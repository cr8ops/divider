package main

import (
	"fmt"

	"github.com/cr8ops/divider/internal/config"
)

func main() {
	cfg, err := config.Read("./data/config.txt")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	fmt.Printf("Video Path: %s\n", cfg.VideoPath)
}
