package playlist

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Track struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Channel   string    `json:"channel"`
	Duration  int       `json:"duration"`
	Thumbnail string    `json:"thumbnail"`
	AddedAt   time.Time `json:"added_at"`
}

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS playlist (
			id        TEXT PRIMARY KEY,
			title     TEXT NOT NULL,
			channel   TEXT NOT NULL,
			duration  INTEGER NOT NULL,
			thumbnail TEXT NOT NULL,
			added_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *Store) Add(t Track) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO playlist (id, title, channel, duration, thumbnail)
		VALUES (?, ?, ?, ?, ?)
	`, t.ID, t.Title, t.Channel, t.Duration, t.Thumbnail)
	return err
}

func (s *Store) List() ([]Track, error) {
	rows, err := s.db.Query(`
		SELECT id, title, channel, duration, thumbnail, added_at
		FROM playlist ORDER BY added_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.Channel, &t.Duration, &t.Thumbnail, &t.AddedAt); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (s *Store) Remove(id string) error {
	_, err := s.db.Exec(`DELETE FROM playlist WHERE id = ?`, id)
	return err
}
