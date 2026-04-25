package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://ws.audioscrobbler.com/2.0/"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type SimilarTrack struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type SimilarArtist struct {
	Name string `json:"name"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SimilarTracks(ctx context.Context, artist, track string, limit int) ([]SimilarTrack, error) {
	params := url.Values{}
	params.Set("method", "track.getSimilar")
	params.Set("artist", artist)
	params.Set("track", track)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	var result struct {
		SimilarTracks struct {
			Track []struct {
				Name   string `json:"name"`
				Artist struct {
					Name string `json:"name"`
				} `json:"artist"`
			} `json:"track"`
		} `json:"similartracks"`
	}

	if err := c.get(ctx, params, &result); err != nil {
		return nil, err
	}

	var tracks []SimilarTrack
	for _, t := range result.SimilarTracks.Track {
		tracks = append(tracks, SimilarTrack{
			Name:   t.Name,
			Artist: t.Artist.Name,
		})
	}
	return tracks, nil
}

func (c *Client) SimilarArtists(ctx context.Context, artist string, limit int) ([]SimilarArtist, error) {
	params := url.Values{}
	params.Set("method", "artist.getSimilar")
	params.Set("artist", artist)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	var result struct {
		SimilarArtists struct {
			Artist []struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"similarartists"`
	}

	if err := c.get(ctx, params, &result); err != nil {
		return nil, err
	}

	var artists []SimilarArtist
	for _, a := range result.SimilarArtists.Artist {
		artists = append(artists, SimilarArtist{Name: a.Name})
	}
	return artists, nil
}

func (c *Client) get(ctx context.Context, params url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("last.fm returned %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
