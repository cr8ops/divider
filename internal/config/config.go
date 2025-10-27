package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	VideoPath string
	Chapters  []Chapter
}

type Chapter struct {
	Start string
	End   string
}

func Read(configPath string) (*Config, error) {
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	config := &Config{
		Chapters: []Chapter{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse video path line
		if strings.HasPrefix(line, "video:") {
			config.VideoPath = strings.TrimSpace(strings.TrimPrefix(line, "video:"))
			continue
		}

		// Parse chapter lines (assuming format: "start-end")
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				chapter := Chapter{
					Start: strings.TrimSpace(parts[0]),
					End:   strings.TrimSpace(parts[1]),
				}
				config.Chapters = append(config.Chapters, chapter)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}

	return config, nil
}
