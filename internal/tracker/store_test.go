package tracker

import (
	"testing"
	"time"
)

func TestRecordAndSignals(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Record a full listen
	if err := store.RecordPlay(PlayEvent{
		VideoID:         "abc123",
		Artist:          "Radiohead",
		Title:           "Creep",
		DurationSeconds: 237,
		ListenSeconds:   237,
		Skipped:         false,
		PlayedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("RecordPlay failed: %v", err)
	}

	// Record a skip
	if err := store.RecordPlay(PlayEvent{
		VideoID:         "def456",
		Artist:          "Radiohead",
		Title:           "No Surprises",
		DurationSeconds: 228,
		ListenSeconds:   30,
		Skipped:         true,
		PlayedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("RecordPlay failed: %v", err)
	}

	signals, err := store.ArtistSignals()
	if err != nil {
		t.Fatalf("ArtistSignals failed: %v", err)
	}

	sig, ok := signals["Radiohead"]
	if !ok {
		t.Fatal("expected Radiohead signal, got none")
	}

	t.Logf("artist=%s plays=%d skips=%d skipRate=%.2f",
		sig.Artist, sig.PlayCount, sig.SkipCount, sig.SkipRate)

	if sig.PlayCount != 2 {
		t.Fatalf("expected 2 plays, got %d", sig.PlayCount)
	}
	if sig.SkipCount != 1 {
		t.Fatalf("expected 1 skip, got %d", sig.SkipCount)
	}
	if sig.SkipRate != 0.5 {
		t.Fatalf("expected skip rate 0.5, got %.2f", sig.SkipRate)
	}
}
