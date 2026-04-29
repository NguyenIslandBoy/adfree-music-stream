package ytdlp

import (
	"context"
	"os"
	"testing"
)

func TestExtractAudioURL(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping extract test in CI — YouTube blocks cloud provider IPs")
	}

	c := NewClient(3)
	url, err := c.ExtractAudioURL(context.Background(), "XFkzRNyygfk")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if url == "" {
		t.Fatal("expected a URL, got empty string")
	}
	t.Logf("audio URL: %s...", url[:60])
}
