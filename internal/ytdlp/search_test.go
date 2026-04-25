package ytdlp

import (
	"context"
	"testing"
)

func TestSearch(t *testing.T) {
	c := NewClient(3)
	results, err := c.Search(context.Background(), "radiohead creep", 3)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}
	for _, r := range results {
		t.Logf("id=%s title=%s channel=%s duration=%ds", r.ID, r.Title, r.Channel, r.Duration)
	}
}
