package tracker

import (
	"testing"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/ytdlp"
)

func TestRerank(t *testing.T) {
	tracks := []ytdlp.TrackMeta{
		{ID: "1", Title: "Track A", Channel: "Radiohead"},
		{ID: "2", Title: "Track B", Channel: "Pixies"},
		{ID: "3", Title: "Track C", Channel: "Muse"},
	}

	// Simulate: user plays Muse a lot, skips Radiohead often
	signals := map[string]*ArtistSignal{
		"Radiohead": {Artist: "Radiohead", PlayCount: 5, SkipCount: 4, SkipRate: 0.8},
		"Muse":      {Artist: "Muse", PlayCount: 10, SkipCount: 1, SkipRate: 0.1},
		"Pixies":    {Artist: "Pixies", PlayCount: 3, SkipCount: 0, SkipRate: 0.0},
	}

	result := Rerank(tracks, signals)

	t.Log("Re-ranked order:")
	for i, r := range result {
		t.Logf("  %d. %s (%s)", i+1, r.Title, r.Channel)
	}

	// Muse should be first — high play count, low skip rate
	if result[0].Channel != "Muse" {
		t.Fatalf("expected Muse first, got %s", result[0].Channel)
	}

	// Radiohead should be last — skip penalty dominates
	if result[len(result)-1].Channel != "Radiohead" {
		t.Fatalf("expected Radiohead last, got %s", result[len(result)-1].Channel)
	}
}
