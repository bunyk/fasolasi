package notes

import (
	"strings"
)

type Pitch struct {
	Frequency float64
	Name      string
	Bottom    float64 // Position relative to first line, in line intervals
	IsHalf    bool    // true if this pitch represents half tone
}

func (p Pitch) Title() string {
	t := strings.ToUpper(p.Name)
	t = strings.ReplaceAll(t, "IS", "#")
	t = strings.ReplaceAll(t, "ES", "b")
	t = strings.ReplaceAll(t, "'", "")
	return t
}

var C = Pitch{523.25, "c", -1, false}

var FluteRange []Pitch
var octave = []Pitch{
	C,
	{554.37, "cis", -0.75, true},
	{587.33, "d", -0.5, false},
	{622.25, "dis", -0.25, true},
	{659.25, "e", 0, false},
	{698.46, "f", 0.5, false},
	{739.99, "fis", 0.75, true},
	{783.99, "g", 1, false},
	{830.61, "gis", 1.25, true},
	{880.00, "a", 1.5, false},
	{932.33, "bes", 1.75, true},
	{987.77, "b", 2, false},
}

var Pause = Pitch{-1.0, "p", 0, false}

var PitchByName map[string]Pitch

func init() {
	FluteRange = append(FluteRange, Pause)
	FluteRange = append(FluteRange, octave...)
	for _, n := range octave {
		FluteRange = append(FluteRange, Pitch{
			Frequency: n.Frequency * 2.0,
			Name:      n.Name + "'",
			Bottom:    n.Bottom + 3.5,
			IsHalf:    n.IsHalf,
		})
	}
	FluteRange = append(FluteRange, Pitch{2093.00, "c''", 6, false})

	PitchByName = make(map[string]Pitch)
	for _, p := range FluteRange {
		PitchByName[p.Name] = p
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
