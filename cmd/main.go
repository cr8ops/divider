/*
Package main is CLI for video divider.
*/
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cr8ops/divider/internal/config"
	"github.com/cr8ops/divider/internal/driver/ffmpeg"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if err := run(l); err != nil {
		l.Error("running", slog.Any("error", err))
	}
}

func run(l *slog.Logger) error {
	cfg, err := config.Read("./data/config.txt")
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	l.Info("video", slog.String("path", cfg.VideoPath))
	mixer := ffmpeg.NewFFmpeg()
	for i, chapter := range cfg.Chapters {
		l.Info("chapter", slog.Int("number", i+1), slog.String("start", chapter.Start), slog.String("end", chapter.End))
		if err := mixer.DivideVideo(l, cfg.VideoPath, fmt.Sprintf("%s/%d.mp4", cfg.OutputPath, i+1), chapter.Start, chapter.End); err != nil {
			l.Error("dividing video", slog.Any("error", err))
		}
	}

	return nil
}
