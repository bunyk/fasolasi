package notes

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type Pitch struct {
	Frequency float64
	Name      string
	Bottom    float64 // Position relative to first line, in line intervals
	Height    float64 // Height (of rectangle representing key) in line intervals
}

var FluteRange []Pitch
var octave = []Pitch{
	{523.25, "c", -1, 0.5},
	{554.37, "cis", -0.75, 0.25},
	{587.33, "d", -0.5, 0.5},
	{622.25, "dis", -0.25, 0.25},
	{659.25, "e", 0, 0.5},
	{698.46, "f", 0.5, 0.5},
	{739.99, "fis", 0.75, 0.25},
	{783.99, "g", 1, 0.5},
	{830.61, "gis", 1.25, 0.25},
	{880.00, "a", 1.5, 0.5},
	{932.33, "bes", 1.75, 0.25},
	{987.77, "b", 2, 0.5},
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
	FluteRange = append(FluteRange, Pitch{2093.00, "c''", 6, 0.5})

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

func NoteFromMatch(parts []string, defaultDuration, fullDuration float64) (SongNote, error) {
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

func ReadSong(filename string) (song []SongNote, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	matches := noteRe.FindAllStringSubmatch(string(data), -1)

	time := 0.0
	fullDuration := 4.0
	defaultDuration := fullDuration / 4
	for _, match := range matches {
		n, err := NoteFromMatch(match, defaultDuration, fullDuration)
		if err != nil {
			return nil, err
		}
		n.Time = time
		time += n.Duration
		defaultDuration = n.Duration
		song = append(song, n)
	}
	return song, nil
}
