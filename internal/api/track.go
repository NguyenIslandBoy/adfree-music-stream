package api

import (
	"encoding/json"
	"net/http"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/tracker"
)

type playRequest struct {
	VideoID         string `json:"video_id"`
	Artist          string `json:"artist"`
	Title           string `json:"title"`
	DurationSeconds int    `json:"duration_seconds"`
	ListenSeconds   int    `json:"listen_seconds"`
	Skipped         bool   `json:"skipped"`
}

func (s *Server) handleTrackPlay(w http.ResponseWriter, r *http.Request) {
	var req playRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.VideoID == "" || req.Artist == "" {
		writeError(w, "video_id and artist are required", http.StatusBadRequest)
		return
	}

	// Auto-detect skip: listened to less than 30% of track
	if req.DurationSeconds > 0 {
		ratio := float64(req.ListenSeconds) / float64(req.DurationSeconds)
		req.Skipped = ratio < 0.3
	}

	err := s.tracker.RecordPlay(tracker.PlayEvent{
		VideoID:         req.VideoID,
		Artist:          req.Artist,
		Title:           req.Title,
		DurationSeconds: req.DurationSeconds,
		ListenSeconds:   req.ListenSeconds,
		Skipped:         req.Skipped,
	})
	if err != nil {
		writeError(w, "failed to record play", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{
		"status":  "recorded",
		"skipped": req.Skipped,
	})
}

func (s *Server) handleTrackPlays(w http.ResponseWriter, r *http.Request) {
	signals, err := s.tracker.ArtistSignals()
	if err != nil {
		writeError(w, "failed to get signals", http.StatusInternalServerError)
		return
	}
	writeJSON(w, signals)
}
