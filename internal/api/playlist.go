package api

import (
	"encoding/json"
	"net/http"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/playlist"
	"github.com/gorilla/mux"
)

func (s *Server) handlePlaylistList(w http.ResponseWriter, r *http.Request) {
	tracks, err := s.playlist.List()
	if err != nil {
		writeError(w, "failed to list playlist", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"tracks": tracks,
		"total":  len(tracks),
	})
}

func (s *Server) handlePlaylistAdd(w http.ResponseWriter, r *http.Request) {
	var track playlist.Track
	if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if track.ID == "" || track.Title == "" {
		writeError(w, "id and title are required", http.StatusBadRequest)
		return
	}
	if err := s.playlist.Add(track); err != nil {
		writeError(w, "failed to add track", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"status": "added", "id": track.ID})
}

func (s *Server) handlePlaylistRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		writeError(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := s.playlist.Remove(id); err != nil {
		writeError(w, "failed to remove track", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "removed", "id": id})
}
