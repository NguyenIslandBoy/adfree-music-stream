package tracker

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type PlayEvent struct {
	VideoID         string
	Artist          string
	Title           string
	DurationSeconds int
	ListenSeconds   int
	Skipped         bool
	PlayedAt        time.Time
}

type ArtistSignal struct {
	Artist    string
	PlayCount int
	SaveCount int
	SkipCount int
	SkipRate  float64
}

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open tracker db: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate tracker: %w", err)
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS plays (
			id               INTEGER PRIMARY KEY AUTOINCREMENT,
			video_id         TEXT NOT NULL,
			artist           TEXT NOT NULL,
			title            TEXT NOT NULL,
			duration_seconds INTEGER NOT NULL,
			listen_seconds   INTEGER NOT NULL,
			skipped          BOOLEAN NOT NULL DEFAULT 0,
			played_at        DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *Store) RecordPlay(e PlayEvent) error {
	_, err := s.db.Exec(`
		INSERT INTO plays (video_id, artist, title, duration_seconds, listen_seconds, skipped, played_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, e.VideoID, e.Artist, e.Title, e.DurationSeconds, e.ListenSeconds, e.Skipped, time.Now())
	return err
}

// ArtistSignals returns play/save/skip counts per artist
// Only considers plays from the last 90 days to keep recommendations fresh
func (s *Store) ArtistSignals() (map[string]*ArtistSignal, error) {
	rows, err := s.db.Query(`
		SELECT
			artist,
			COUNT(*) as play_count,
			SUM(CASE WHEN skipped = 1 THEN 1 ELSE 0 END) as skip_count
		FROM plays
		WHERE played_at >= datetime('now', '-90 days')
		GROUP BY artist
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	signals := map[string]*ArtistSignal{}
	for rows.Next() {
		var a ArtistSignal
		if err := rows.Scan(&a.Artist, &a.PlayCount, &a.SkipCount); err != nil {
			return nil, err
		}
		if a.PlayCount > 0 {
			a.SkipRate = float64(a.SkipCount) / float64(a.PlayCount)
		}
		signals[a.Artist] = &a
	}
	return signals, nil
}

// InjectSaveCounts adds playlist save counts into existing signals
func (s *Store) InjectSaveCounts(signals map[string]*ArtistSignal, savedTracks []string) {
	// savedTracks is a list of artist names from the playlist
	for _, artist := range savedTracks {
		if sig, ok := signals[artist]; ok {
			sig.SaveCount++
		} else {
			signals[artist] = &ArtistSignal{Artist: artist, SaveCount: 1}
		}
	}
}
