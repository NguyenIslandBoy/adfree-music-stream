package tracker

import (
	"github.com/NguyenIslandBoy/adfree-music-stream/internal/ytdlp"
)

const (
	weightPlayCount = 0.3
	weightSaveCount = 0.8
	weightSkipRate  = 1.5
)

type ScoredTrack struct {
	Track ytdlp.TrackMeta
	Score float64
}

// Rerank takes a list of recommended tracks and re-ranks them
// based on user play/save/skip signals for each artist
func Rerank(tracks []ytdlp.TrackMeta, signals map[string]*ArtistSignal) []ytdlp.TrackMeta {
	if len(signals) == 0 {
		return tracks // no data yet, return Last.fm order as-is
	}

	scored := make([]ScoredTrack, len(tracks))
	for i, t := range tracks {
		// Base score — earlier in Last.fm results = higher base score
		base := float64(len(tracks)-i) / float64(len(tracks))

		score := base
		if sig, ok := signals[t.Channel]; ok {
			score += float64(sig.PlayCount) * weightPlayCount
			score += float64(sig.SaveCount) * weightSaveCount
			score -= sig.SkipRate * weightSkipRate
		}

		scored[i] = ScoredTrack{Track: t, Score: score}
	}

	// Sort descending by score
	sortScoredTracks(scored)

	result := make([]ytdlp.TrackMeta, len(scored))
	for i, s := range scored {
		result[i] = s.Track
	}
	return result
}

// simple insertion sort — small slices only (≤20 tracks)
func sortScoredTracks(tracks []ScoredTrack) {
	for i := 1; i < len(tracks); i++ {
		key := tracks[i]
		j := i - 1
		for j >= 0 && tracks[j].Score < key.Score {
			tracks[j+1] = tracks[j]
			j--
		}
		tracks[j+1] = key
	}
}
