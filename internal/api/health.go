package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	version, err := s.ytdlp.Version(r.Context())
	if err != nil {
		version = "unavailable"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":        "ok",
		"ytdlp_version": version,
		"cache_items":   s.cache.ItemCount(),
	})
}
