package api

import (
	"net/http"
	"time"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/cache"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/lastfm"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/playlist"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/tracker"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/ytdlp"
	"github.com/gorilla/mux"
)

type Server struct {
	ytdlp    *ytdlp.Client
	cache    *cache.Cache
	lastfm   *lastfm.Client
	playlist *playlist.Store
	tracker  *tracker.Store
}

func NewRouter(ytClient *ytdlp.Client, c *cache.Cache, lfClient *lastfm.Client, pl *playlist.Store, tr *tracker.Store) *mux.Router {
	s := &Server{
		ytdlp:    ytClient,
		cache:    c,
		lastfm:   lfClient,
		playlist: pl,
		tracker:  tr,
	}

	r := mux.NewRouter()

	// Frontend — no auth
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	r.Handle("/", http.FileServer(http.Dir("./static")))

	// API — auth required
	api := r.PathPrefix("").Subrouter()
	api.Use(authMiddleware)
	api.HandleFunc("/health", s.handleHealth).Methods("GET")
	api.HandleFunc("/search", s.handleSearch).Methods("GET")
	api.HandleFunc("/stream/{id}", s.handleStream).Methods("GET")
	api.HandleFunc("/recommendations/{id}", s.handleRecommendations).Methods("GET")
	api.HandleFunc("/playlist", s.handlePlaylistList).Methods("GET")
	api.HandleFunc("/playlist", s.handlePlaylistAdd).Methods("POST")
	api.HandleFunc("/playlist/{id}", s.handlePlaylistRemove).Methods("DELETE")
	api.HandleFunc("/track/play", s.handleTrackPlay).Methods("POST")
	api.HandleFunc("/track/plays", s.handleTrackPlays).Methods("GET")

	return r
}

const searchTTL = 10 * time.Minute
const streamTTL = 5 * time.Hour
