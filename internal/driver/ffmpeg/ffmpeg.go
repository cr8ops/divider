/*
Package ffmpeg is a general wrapper around ffmpeg binary
*/
package ffmpeg

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// BinaryName is a ffmpeg binary name
const BinaryName = "ffmpeg"

// DefaultArgs are default arguments used in every ffmpeg binary execution
var DefaultArgs = []string{
	"-hide_banner",
	"-loglevel", "error",
	"-y",
}

// FFmpeg represents FFmpeg wrapper
type FFmpeg struct {
}

// NewFFmpeg constructs the FFmpeg wrapper
func NewFFmpeg() *FFmpeg {
	return &FFmpeg{}
}

// Run runs ffmpeg command
func run(l *slog.Logger, cmd *exec.Cmd) error {
	var out, outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr

	lr := l.With("cmd", cmd)
	lr.Info("running a command")
	err := cmd.Run()
	strOutErr := outErr.String()
	lr.Info("finished a command",
		slog.String("stdout", out.String()),
		slog.String("stderr", strOutErr),
	)
	// Ignore error "inadequate AVFrame plane padding", beause
	// it doesn't seem to block concatenating audio files
	// and it appears randomly
	if err != nil && !strings.Contains(err.Error(), "exit status 234") {
		return fmt.Errorf("running a command: %w", err)
	}
	if len(strOutErr) > 0 && strings.Contains(strings.ToLower(strOutErr), "error") {
		return fmt.Errorf("running a command: %s", strOutErr)
	}
	return nil
}

// DivideVideo divides video into multiple parts
func (ff *FFmpeg) DivideVideo(l *slog.Logger, videoFilePath, dstFilePath, start, end string) error {
	args := DefaultArgs
	// ffmpeg -i input.mp4 -ss 00:00:00 -to 00:10:00 -c copy output_part1.mp4
	args = append(args,
		"-i", videoFilePath,
		"-ss", start,
		"-to", end,
		"-copyts",
		dstFilePath,
	)

	if err := run(l, exec.Command(BinaryName, args...)); err != nil { //nolint:gosec
		return fmt.Errorf("executting a command: %w", err)
	}

	return nil
}

/*
func (ff *FFmpeg) DivideVideo(l *slog.Logger, videoFilePath, dstFilePath, start, end string) error {
	args := DefaultArgs
	// ffmpeg -i input.mp4 -ss 00:00:00 -to 00:10:00 -c copy output_part1.mp4
	args = append(args,
		"-i", videoFilePath,
		"-ss", start,
		"-t", end,
		dstFilePath,
	)

	if err := run(l, exec.Command(BinaryName, args...)); err != nil { //nolint:gosec
		return fmt.Errorf("executting a command: %w", err)
	}

	return nil
}*/
