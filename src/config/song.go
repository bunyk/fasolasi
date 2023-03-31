package config

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bunyk/fasolasi/src/notes"
)

type Song struct {
	Name  string `yaml:"name"`
	Notes string `yaml:"notes"`
}

var noteRe = regexp.MustCompile(`([a-z']+)(\d+)?(.?)`)

func noteFromMatch(parts []string, defaultDuration, fullDuration float64) (notes.SongNote, error) {
	p := notes.PitchByName[parts[1]]
	note := notes.SongNote{}
	if p.Name == "" {
		return note, fmt.Errorf("Unknown note: %#v", parts[1])
	}
	duration := defaultDuration
	if parts[2] != "" {
		nd, err := strconv.Atoi(parts[2])
		if err != nil {
			return note, fmt.Errorf("Failed to parse duration %s: %w", parts[2], err)
		}
		duration = fullDuration / float64(nd)
	}
	if parts[3] == "." {
		duration *= 1.5
	}
	return notes.SongNote{
		Pitch:    p,
		Duration: duration,
	}, nil
}

// bpm - beats per minute
// beat is a note in denominator of the time signature. Ex: in 4/4, 3/4 - beat is quarter note
// So duration of full note for /4 tempo is 60 / bpm * 4 = 240 / bpm. Ex, for 60 bpm - 4 seconds. 120 bpm - 2 seconds.
func (s Song) ParseNotes(fullDuration float64) (song []notes.SongNote, err error) {
	matches := noteRe.FindAllStringSubmatch(s.Notes, -1)
	time := TimeBeforeFirstNote // give some initial time to prepare for first note
	defaultDuration := fullDuration / 4
	for _, match := range matches {
		n, err := noteFromMatch(match, defaultDuration, fullDuration)
		if err != nil {
			return nil, err
		}
		n.Time = time
		time += n.Duration + BreathInterval*fullDuration
		defaultDuration = n.Duration
		song = append(song, n)
	}
	return song, nil
}
