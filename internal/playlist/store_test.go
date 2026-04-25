package playlist

import (
	"testing"
)

func TestStore(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	track := Track{
		ID:        "XFkzRNyygfk",
		Title:     "Radiohead - Creep",
		Channel:   "Radiohead",
		Duration:  237,
		Thumbnail: "https://i.ytimg.com/vi/XFkzRNyygfk/hqdefault.jpg",
	}

	// Add
	if err := store.Add(track); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// List
	tracks, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(tracks))
	}
	t.Logf("track: %s - %s", tracks[0].ID, tracks[0].Title)

	// Remove
	if err := store.Remove(track.ID); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	tracks, _ = store.List()
	if len(tracks) != 0 {
		t.Fatal("expected empty playlist after remove")
	}
	t.Log("remove ok")
}
