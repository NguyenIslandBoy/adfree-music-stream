package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/tracker"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/ytdlp"
	"github.com/gorilla/mux"
)

func (s *Server) handleRecommendations(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		writeError(w, "missing video id", http.StatusBadRequest)
		return
	}

	cacheKey := "recommendations:" + id
	if cached, ok := s.cache.Get(cacheKey); ok {
		writeJSON(w, cached)
		return
	}

	track, err := s.resolveTrack(r.Context(), id)
	if err != nil {
		writeError(w, "could not resolve track metadata", http.StatusInternalServerError)
		return
	}

	artist, title := parseArtistTitle(track.Title, track.Channel)
	seen := map[string]bool{id: true}
	var results []ytdlp.TrackMeta

	// Stage 1: Last.fm similar tracks
	similarTracks, err := s.lastfm.SimilarTracks(r.Context(), artist, title, 10)
	if err == nil && len(similarTracks) > 0 {
		for _, st := range similarTracks {
			query := fmt.Sprintf("%s %s", st.Artist, st.Name)
			tracks, err := s.ytdlp.Search(r.Context(), query, 1)
			if err != nil || len(tracks) == 0 {
				continue
			}
			t := tracks[0]
			if seen[t.ID] {
				continue
			}
			seen[t.ID] = true
			if !isValidTrack(t) {
				continue
			}
			results = append(results, t)
		}
	}

	// Stage 2: Last.fm similar artists fallback
	if len(results) < 5 {
		similarArtists, err := s.lastfm.SimilarArtists(r.Context(), artist, 10)
		if err == nil && len(similarArtists) > 0 {
			for _, sa := range similarArtists {
				tracks, err := s.ytdlp.Search(r.Context(), sa.Name, 1)
				if err != nil || len(tracks) == 0 {
					continue
				}
				t := tracks[0]
				if seen[t.ID] {
					continue
				}
				seen[t.ID] = true
				if !isValidTrack(t) {
					continue
				}
				results = append(results, t)
			}
		}
	}

	// Stage 3: yt-dlp keyword fallback for unknown tracks (e.g. Vietnamese/TikTok music)
	if len(results) < 5 {
		queries := buildFallbackQueries(artist, title, track.Channel)
		for _, q := range queries {
			tracks, err := s.ytdlp.Search(r.Context(), q, 2)
			if err != nil {
				continue
			}
			for _, t := range tracks {
				if seen[t.ID] {
					continue
				}
				seen[t.ID] = true
				if !isValidTrack(t) {
					continue
				}
				results = append(results, t)
			}
			if len(results) >= 8 {
				break
			}
		}
	}

	// Rerank based on user signals
	signals, err := s.tracker.ArtistSignals()
	if err == nil && len(signals) > 0 {
		// Inject playlist save counts as signals
		savedTracks, err := s.playlist.List()
		if err == nil {
			artistNames := make([]string, len(savedTracks))
			for i, t := range savedTracks {
				artistNames[i] = t.Channel
			}
			s.tracker.InjectSaveCounts(signals, artistNames)
		}
		results = tracker.Rerank(results, signals)
	}

	resp := map[string]any{
		"recommendations": results,
		"total":           len(results),
	}
	s.cache.Set(cacheKey, resp, 2*time.Minute)
	writeJSON(w, resp)
}

// buildFallbackQueries generates YouTube search queries from track metadata
// used when Last.fm has no data (e.g. Vietnamese, TikTok, niche music)
func buildFallbackQueries(artist, title, channel string) []string {
	var queries []string

	// Search by artist name directly
	if artist != channel {
		queries = append(queries, artist)
	}

	// Search by channel
	queries = append(queries, channel)

	// Extract meaningful keywords from title (strip common noise words)
	keywords := extractKeywords(title)
	if keywords != "" {
		queries = append(queries, keywords)
	}

	return queries
}

func extractKeywords(title string) string {
	// Strip common noise patterns from titles
	noise := []string{
		"(Official Music Video)", "(Official Video)", "(Lyrics)",
		"(Audio)", "(MV)", "Chapter 1", "Chapter 2",
		"[Vietsub]", "[Pinyin]", "Remix", "Ver",
	}
	result := title
	for _, n := range noise {
		result = strings.ReplaceAll(result, n, "")
	}
	result = strings.TrimSpace(result)

	// Take first 4 words as keywords
	words := strings.Fields(result)
	if len(words) > 4 {
		words = words[:4]
	}
	return strings.Join(words, " ")
}

// parseArtistTitle tries to split "Artist - Title" format
// falls back to channel as artist if no dash found
func parseArtistTitle(title, channel string) (artist, track string) {
	if parts := strings.SplitN(title, " - ", 2); len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return channel, title
}

func isValidTrack(t ytdlp.TrackMeta) bool {
	return t.Duration > 30 && t.Duration < 3600 // ignore channels (0s) and hour-long compilations
}
