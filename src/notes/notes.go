package notes

import (
	"fmt"
	"io"
	"os"
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

var pitchByName map[string]Pitch

func init() {
	FluteRange = append(FluteRange, Pitch{-1.0, "p", 0, 0}) // pause
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
	return sn.Time + sn.Duration
}

func NoteFromString(s string) (SongNote, error) {
	p := pitchByName[s] // TODO: parse durations
	if p.Name == "" {
		return SongNote{}, fmt.Errorf("Unknown note: %#v", s)
	}
	return SongNote{
		Pitch:    p,
		Duration: 1.0,
	}, nil
}

func ReadSong(filename string) (song []SongNote, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var token string
	time := 0.0
	for {
		_, err := fmt.Fscanf(f, "%s", &token)
		if err != nil {
			if err == io.EOF {
				return song, nil
			}
			return nil, err
		}
		n, err := NoteFromString(token)
		if err != nil {
			return nil, err
		}
		n.Time = time
		time += n.Duration
		song = append(song, n)
	}
	return nil, nil
}
