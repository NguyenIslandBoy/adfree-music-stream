package lastfm

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestSimilarTracks(t *testing.T) {
	godotenv.Load("../../.env")
	key := os.Getenv("LASTFM_API_KEY")
	if key == "" {
		t.Skip("LASTFM_API_KEY not set")
	}

	c := NewClient(key)
	tracks, err := c.SimilarTracks(context.Background(), "Radiohead", "Creep", 5)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(tracks) == 0 {
		t.Fatal("expected results, got none")
	}
	for _, tr := range tracks {
		t.Logf("similar track: %s - %s", tr.Artist, tr.Name)
	}
}

func TestSimilarArtists(t *testing.T) {
	godotenv.Load("../../.env")
	key := os.Getenv("LASTFM_API_KEY")
	if key == "" {
		t.Skip("LASTFM_API_KEY not set")
	}

	c := NewClient(key)
	artists, err := c.SimilarArtists(context.Background(), "Radiohead", 5)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	for _, a := range artists {
		t.Logf("similar artist: %s", a.Name)
	}
}
