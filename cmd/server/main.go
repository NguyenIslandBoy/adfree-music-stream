package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NguyenIslandBoy/adfree-music-stream/internal/api"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/cache"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/lastfm"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/playlist"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/tracker"
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/ytdlp"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	allowedOrigin := getEnv("ALLOWED_ORIGIN", "http://localhost:8080")
	maxWorkers := getEnvInt("MAX_CONCURRENT_YTDLP", 3)

	ytClient := ytdlp.NewClient(maxWorkers)
	urlCache := cache.New(5*time.Hour, 10*time.Minute)
	lfClient := lastfm.NewClient(mustGetEnv("LASTFM_API_KEY"))
	pl, err := playlist.NewStore(getEnv("DB_PATH", "playlist.db"))
	tr, err := tracker.NewStore(getEnv("DB_PATH", "playlist.db"))
	if err != nil {
		log.Fatalf("failed to open tracker db: %v", err)
	}
	router := api.NewRouter(ytClient, urlCache, lfClient, pl, tr)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigin},
		AllowedMethods: []string{"GET"},
	})
	handler := c.Handler(router)

	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}
