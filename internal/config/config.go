/*
Package config provides utilities for managing configuration files.
*/
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	chapterPartsDelimiter = "-"
	chapterPartsNumber    = 2
)

// Config represents the configuration for video division.
type Config struct {
	VideoPath string
	Chapters  []Chapter
}

// Chapter represents a video chapter with start and end times.
type Chapter struct {
	Start string
	End   string
}

// Read reads the configuration from the specified file path.
func Read(configPath string) (cfg *Config, err error) {
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			err = fmt.Errorf("closing file: %w", err)
		}
	}()

	cfg = &Config{
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
			cfg.VideoPath = strings.TrimSpace(strings.TrimPrefix(line, "video:"))
			continue
		}

		// Parse chapter lines (assuming format: "start-end")
		if strings.Contains(line, chapterPartsDelimiter) {
			parts := strings.Split(line, chapterPartsDelimiter)
			if len(parts) == chapterPartsNumber {
				chapter := Chapter{
					Start: strings.TrimSpace(parts[0]),
					End:   strings.TrimSpace(parts[1]),
				}
				cfg.Chapters = append(cfg.Chapters, chapter)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}

	return cfg, nil
}
