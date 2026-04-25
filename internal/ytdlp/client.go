package ytdlp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Client struct {
	binary string
	sem    chan struct{}
}

func NewClient(maxConcurrent int) *Client {
	return &Client{
		binary: getEnv("YTDLP_BINARY", "yt-dlp"),
		sem:    make(chan struct{}, maxConcurrent),
	}
}

func (c *Client) acquire(ctx context.Context) error {
	select {
	case c.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(15 * time.Second):
		return fmt.Errorf("yt-dlp pool exhausted, try again later")
	}
}

func (c *Client) release() {
	<-c.sem
}

func (c *Client) Version(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, c.binary, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("yt-dlp not found or not working: %w", err)
	}
	return string(out), nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
