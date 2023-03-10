package notes

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bunyk/fasolasi/src/config"
)

type Pitch struct {
	Frequency float64
	Name      string
	Bottom    float64 // Position relative to first line, in line intervals
	Height    float64 // Height (of rectangle representing key) in line intervals
}

func (p Pitch) Title() string {
	t := strings.ToUpper(p.Name)
	t = strings.ReplaceAll(t, "IS", "#")
	t = strings.ReplaceAll(t, "ES", "b")
	t = strings.ReplaceAll(t, "'", "")
	return t
}

var C = Pitch{523.25, "c", -1, config.WhiteNoteWidth}

var FluteRange []Pitch
var octave = []Pitch{
	C,
	{554.37, "cis", -0.75, config.BlackNoteWidth},
	{587.33, "d", -0.5, config.WhiteNoteWidth},
	{622.25, "dis", -0.25, config.BlackNoteWidth},
	{659.25, "e", 0, config.WhiteNoteWidth},
	{698.46, "f", 0.5, config.WhiteNoteWidth},
	{739.99, "fis", 0.75, config.BlackNoteWidth},
	{783.99, "g", 1, config.WhiteNoteWidth},
	{830.61, "gis", 1.25, config.BlackNoteWidth},
	{880.00, "a", 1.5, config.WhiteNoteWidth},
	{932.33, "bes", 1.75, config.BlackNoteWidth},
	{987.77, "b", 2, config.WhiteNoteWidth},
}

var Pause = Pitch{-1.0, "p", 0, 0}

var pitchByName map[string]Pitch

func init() {
	FluteRange = append(FluteRange, Pause)
	FluteRange = append(FluteRange, octave...)
	for _, n := range octave {
		FluteRange = append(FluteRange, Pitch{
			Frequency: n.Frequency * 2.0,
			Name:      n.Name + "'",
			Bottom:    n.Bottom + 3.5,
			Height:    n.Height,
		})
	}
	FluteRange = append(FluteRange, Pitch{2093.00, "c''", 6, config.WhiteNoteWidth})

	pitchByName = make(map[string]Pitch)
	for _, p := range FluteRange {
		pitchByName[p.Name] = p
	}
}

func (p Pitch) HasAdditionalLine() bool {
	if 0 <= p.Bottom && p.Bottom <= 4 {
		return false // On existing lines
	}
	return int(p.Bottom*4)%4 == 0
}

func GuessNote(frequency float64) (Pitch, int) {
	min := 0
	max := len(FluteRange) - 1
	for {
		if frequency <= FluteRange[min].Frequency {
			return FluteRange[min], min
		}
		if frequency >= FluteRange[max].Frequency {
			return FluteRange[max], max
		}
		if max-min <= 1 {
			toMax := FluteRange[max].Frequency - frequency
			toMin := frequency - FluteRange[min].Frequency
			if toMax < toMin {
				return FluteRange[max], max
			}
			return FluteRange[min], min
		}
		middle := (min + max) / 2
		if frequency <= FluteRange[middle].Frequency {
			max = middle
		} else {
			min = middle
		}
	}
}

type SongNote struct {
	Pitch    Pitch
	Time     float64
	Duration float64
}

func (sn SongNote) End() float64 {
	if sn.Duration < 0 {
		return -1.0 // Endless :)
	}
	return sn.Time + sn.Duration
}

var noteRe = regexp.MustCompile(`([a-z']+)(\d+)?(.?)`)

func noteFromMatch(parts []string, defaultDuration, fullDuration float64) (SongNote, error) {
	p := pitchByName[parts[1]]
	note := SongNote{}
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
	return SongNote{
		Pitch:    p,
		Duration: duration,
	}, nil
}

// bpm - beats per minute
// beat is a note in denominator of the time signature. Ex: in 4/4, 3/4 - beat is quarter note
// So duration of full note for /4 tempo is 60 / bpm * 4 = 240 / bpm. Ex, for 60 bpm - 4 seconds. 120 bpm - 2 seconds.
func ReadSong(filename string, fullDuration float64) (song []SongNote, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	matches := noteRe.FindAllStringSubmatch(string(data), -1)

	time := config.TimeBeforeFirstNote // give some initial time to prepare for first note
	defaultDuration := fullDuration / 4
	for _, match := range matches {
		n, err := noteFromMatch(match, defaultDuration, fullDuration)
		if err != nil {
			return nil, err
		}
		n.Time = time
		time += n.Duration + config.BreathInterval*fullDuration
		defaultDuration = n.Duration
		song = append(song, n)
	}
	return song, nil
}
