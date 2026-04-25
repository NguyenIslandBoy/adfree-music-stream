package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		writeError(w, "missing video id", http.StatusBadRequest)
		return
	}

	cacheKey := "stream:" + id

	if cached, ok := s.cache.Get(cacheKey); ok {
		http.Redirect(w, r, cached.(string), http.StatusFound)
		return
	}

	url, err := s.ytdlp.ExtractAudioURL(r.Context(), id)
	if err != nil {
		writeError(w, "failed to extract stream", http.StatusInternalServerError)
		return
	}

	s.cache.Set(cacheKey, url, streamTTL)
	http.Redirect(w, r, url, http.StatusFound)
}
