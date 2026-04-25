package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, "missing param: q", http.StatusBadRequest)
		return
	}

	limit := parseIntParam(r, "limit", 10, 20)
	cacheKey := fmt.Sprintf("search:%s:%d", q, limit)

	if cached, ok := s.cache.Get(cacheKey); ok {
		writeJSON(w, cached)
		return
	}

	results, err := s.ytdlp.Search(r.Context(), q, limit)
	if err != nil {
		writeError(w, "search failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"results": results,
		"total":   len(results),
	}
	s.cache.Set(cacheKey, resp, searchTTL)
	writeJSON(w, resp)
}

func parseIntParam(r *http.Request, key string, defaultVal, max int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return defaultVal
	}
	if n > max {
		return max
	}
	return n
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
