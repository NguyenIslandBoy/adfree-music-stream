package ytdlp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func (c *Client) ExtractAudioURL(ctx context.Context, videoID string) (string, error) {
	if err := c.acquire(ctx); err != nil {
		return "", err
	}
	defer c.release()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	args := []string{
		"--get-url",
		"--no-warnings",
		"--no-playlist",
		"-f", "bestaudio",
		"--", videoID,
	}

	out, err := exec.CommandContext(ctx, c.binary, args...).Output()
	if err != nil {
		return "", fmt.Errorf("yt-dlp extract: %w", err)
	}

	url := strings.TrimSpace(string(out))
	if url == "" {
		return "", fmt.Errorf("yt-dlp returned empty URL for id: %s", videoID)
	}

	return url, nil
}
