package ytdlp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type TrackMeta struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Channel    string `json:"channel"`
	Duration   int    `json:"duration"`
	Thumbnail  string `json:"thumbnail"`
	UploadDate string `json:"upload_date"`
}

type flatEntry struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Channel    string  `json:"channel"`
	Duration   float64 `json:"duration"` // yt-dlp returns float, not int
	Thumbnails []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
	UploadDate string `json:"upload_date"`
}

func (c *Client) Search(ctx context.Context, query string, limit int) ([]TrackMeta, error) {
	if err := c.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.release()

	args := []string{
		"--flat-playlist",
		"--dump-json",
		"--no-warnings",
		"--no-playlist",
		"-f", "bestaudio/best",
		"--cookies", getEnv("YTDLP_COOKIES", "cookies.txt"),
		fmt.Sprintf("ytsearch%d:%s", limit, query),
	}

	out, err := exec.CommandContext(ctx, c.binary, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp search: %w", err)
	}

	return parseEntries(out, limit)
}

func parseEntries(out []byte, limit int) ([]TrackMeta, error) {
	var results []TrackMeta

	for _, line := range bytes.Split(out, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		var entry flatEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue // skip malformed lines
		}

		track := TrackMeta{
			ID:         entry.ID,
			Title:      entry.Title,
			Channel:    entry.Channel,
			Duration:   int(entry.Duration),
			UploadDate: entry.UploadDate,
		}

		if len(entry.Thumbnails) > 0 {
			track.Thumbnail = entry.Thumbnails[len(entry.Thumbnails)-1].URL
		}

		results = append(results, track)
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}
