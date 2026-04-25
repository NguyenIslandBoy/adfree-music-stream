# Ad-Free Music Stream

A self-hosted, ad-free music streaming tool built in Go.
Searches YouTube and streams audio directly from YouTube's CDN - no ads, no storage, no third-party player.

---

## Features

- 🔍 YouTube search via `yt-dlp`
- 🎵 Audio-only streaming (no video, no ads)
- 💾 Save tracks to a personal playlist (SQLite)
- 🎯 Recommendations via Last.fm API, re-ranked by your listening behaviour
- ⏭ Autoplay - next recommendation, then playlist fallback
- 📊 Play tracking - skip detection, listen time, per-artist signals
- 🔒 Bearer token authentication
- 💰 ~$0/month to run (Oracle Free Tier or any cheap VPS)

---

## Architecture

```
Browser (HTML audio player)
        │
        │  HTTP REST
        ▼
┌─────────────────────────────────┐
│        Go REST API              │
│  (single binary, your VPS)      │
│                                 │
│  /search    → yt-dlp search     │
│  /stream    → 302 CDN redirect  │
│  /recommendations → Last.fm     │
│                  + yt-dlp       │
│                  + tracker      │
│  /playlist  → SQLite CRUD       │
│  /track/play → play events      │
└─────────────────────────────────┘
        │
        │ 302 redirect
        ▼
  YouTube CDN (googlevideo.com)
  ← browser streams audio here
  ← YouTube pays the bandwidth
```

The Go API never touches the audio bytes. It extracts the CDN URL and returns a `302 redirect` - the browser connects directly to `googlevideo.com`. Your VPS only handles lightweight JSON and redirects.

---

## Tech Stack

| Layer | Choice | Reason |
|---|---|---|
| Language | Go 1.26+ | Fast, single binary, low memory |
| HTTP router | `gorilla/mux` | Simple, well-tested |
| Audio extraction | `yt-dlp -f bestaudio` | Best audio-only extraction, no ads |
| Cache | `patrickmn/go-cache` | In-memory TTL, zero infrastructure |
| Database | `modernc.org/sqlite` | Pure Go SQLite, no C compiler needed |
| Recommendations | Last.fm API + yt-dlp fallback | Rich music graph, works offline for niche music |
| CORS | `rs/cors` | One-line setup for browser access |
| Config | `joho/godotenv` | `.env` file support |
| Concurrency | Semaphore (buffered channel) | Cap parallel yt-dlp processes |
| Container | Docker + Python Alpine | yt-dlp needs Python at runtime |

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Server status + yt-dlp version |
| GET | `/search?q=&limit=` | Search YouTube for tracks |
| GET | `/stream/:id` | 302 redirect to YouTube CDN audio URL |
| GET | `/recommendations/:id` | Last.fm + tracker re-ranked suggestions |
| GET | `/playlist` | List saved tracks |
| POST | `/playlist` | Save a track to playlist |
| DELETE | `/playlist/:id` | Remove a track from playlist |
| POST | `/track/play` | Record a play event |
| GET | `/track/plays` | Artist signals for debugging |

All endpoints except `/` require:
- `Authorization: Bearer <token>` header, or
- `?token=<token>` query param (used for audio streaming redirects)

---

## Getting Started

### Prerequisites

- Go 1.26+
- Python 3.x + `yt-dlp`
- Last.fm API key - [get one free here](https://www.last.fm/api/account/create)

### Install

```bash
git clone https://github.com/NguyenIslandBoy/adfree-music-stream
cd adfree-music-stream
go mod download
pip install yt-dlp
```

### Configure

```bash
cp .env.example .env
```

Edit `.env`:

```env
PORT=8081
ALLOWED_ORIGIN=http://localhost:8081
MAX_CONCURRENT_YTDLP=3
API_TOKEN=your_secret_token_here
LASTFM_API_KEY=your_lastfm_key_here
DB_PATH=playlist.db
```

### Run

```bash
go build ./cmd/server/
./server
```

Open `http://localhost:8081` in your browser.

### Run with Docker

```bash
docker compose up --build
```

---

## Deployment

### Recommended: Hetzner CAX11 (€3.29/month)

ARM instance, 2 vCPU, 4GB RAM - more than enough for personal use.

```bash
# On your VPS
git clone https://github.com/NguyenIslandBoy/adfree-music-stream
cd adfree-music-stream
cp .env.example .env
# edit .env with production values
docker compose up -d
```

### yt-dlp Auto-Update

The Docker image includes a weekly cron job (every Monday at 3am) to keep yt-dlp updated automatically:

```
0 3 * * 1 pip install -U yt-dlp --quiet
```

To update manually on your VPS:

```bash
pip install -U yt-dlp
docker compose restart
```

---

## Project Structure

```
adfree-music-stream/
├── cmd/
│   └── server/
│       └── main.go           # Entry point, config, server bootstrap
├── internal/
│   ├── api/
│   │   ├── router.go         # Route definitions + middleware
│   │   ├── middleware.go     # Bearer token auth
│   │   ├── health.go         # GET /health
│   │   ├── search.go         # GET /search
│   │   ├── stream.go         # GET /stream/:id
│   │   ├── recommendations.go # GET /recommendations/:id
│   │   ├── playlist.go       # GET/POST/DELETE /playlist
│   │   ├── track.go          # POST /track/play
│   │   └── resolve.go        # Track metadata helper
│   ├── cache/
│   │   └── cache.go          # In-memory TTL cache wrapper
│   ├── lastfm/
│   │   └── client.go         # Last.fm API client
│   ├── playlist/
│   │   └── store.go          # SQLite playlist CRUD
│   ├── tracker/
│   │   ├── store.go          # SQLite play events store
│   │   └── scorer.go         # Recommendation re-ranking
│   └── ytdlp/
│       ├── client.go         # yt-dlp subprocess wrapper + semaphore
│       ├── search.go         # Search logic + result parsing
│       └── extract.go        # Audio URL extraction
├── static/
│   └── index.html            # Frontend - search, player, playlist, recommendations
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

---

## How Recommendations Work

```
1. GET /recommendations/:id
        │
        ├─ Stage 1: Last.fm track.getSimilar
        │   → search YouTube for each similar track
        │
        ├─ Stage 2: Last.fm artist.getSimilar (fallback if < 5 results)
        │   → search YouTube for each similar artist
        │
        ├─ Stage 3: yt-dlp keyword fallback (for niche/regional music)
        │   → extract keywords from title, search YouTube directly
        │
        └─ Rerank by user signals:
            score = lastfm_position
                  + (artist_play_count × 0.3)
                  + (artist_save_count × 0.8)
                  - (artist_skip_rate  × 1.5)
```

Play signals are collected automatically - skip detection triggers when you listen to less than 30% of a track's duration.

---

## Limitations

| Concern | Detail |
|---|---|
| yt-dlp breakage | YouTube changes internals monthly - weekly auto-update handles this |
| CDN URL expiry | YouTube CDN URLs valid ~6h - in-memory cache with 5h TTL handles this |
| IP rate limiting | Heavy searching may get throttled - personal use is fine |
| YouTube ToS | Bypassing ads violates ToS - self-hosted personal use is low risk |
| Last.fm coverage | Limited data for niche/regional music - yt-dlp keyword fallback handles this |
| No offline | Streams only, no download or save feature |

---

## Cost

| Item | Cost |
|---|---|
| Hetzner CAX11 VPS | €3.29/month |
| Last.fm API | $0 (free) |
| yt-dlp | $0 (open source) |
| SQLite | $0 (embedded) |
| Bandwidth | $0 (YouTube CDN serves audio) |

---

## License

MIT