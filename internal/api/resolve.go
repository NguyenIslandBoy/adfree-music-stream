package api

import (
	"context"
	"fmt"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/ytdlp"
)

func (s *Server) resolveTrack(ctx context.Context, id string) (*ytdlp.TrackMeta, error) {
	cacheKey := "meta:" + id
	if cached, ok := s.cache.Get(cacheKey); ok {
		if t, ok := cached.(*ytdlp.TrackMeta); ok {
			return t, nil
		}
	}

	// Fall back to searching by ID
	results, err := s.ytdlp.Search(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%s", id), 1)
	if err != nil || len(results) == 0 {
		return nil, fmt.Errorf("could not resolve track: %s", id)
	}

	track := &results[0]
	s.cache.Set(cacheKey, track, searchTTL)
	return track, nil
}
