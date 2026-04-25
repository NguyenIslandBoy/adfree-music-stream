package ytdlp

import (
	"context"
	"testing"
)

func TestVersion(t *testing.T) {
	c := NewClient(3)
	version, err := c.Version(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("yt-dlp version: %s", version)
}
