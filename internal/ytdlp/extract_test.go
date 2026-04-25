package ytdlp

import (
	"context"
	"testing"
)

func TestExtractAudioURL(t *testing.T) {
	c := NewClient(3)
	// Using the Radiohead Creep ID from our search test
	url, err := c.ExtractAudioURL(context.Background(), "XFkzRNyygfk")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if url == "" {
		t.Fatal("expected a URL, got empty string")
	}
	t.Logf("audio URL: %s...", url[:60])
}
