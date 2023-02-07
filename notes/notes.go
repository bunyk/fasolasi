package yin

type note struct {
	Frequency float64
	Name      string
}

var notes = []note{
	{-1.0, "pause"},
	{220.00, "0"},
	{233.08, "1"},
	{246.94, "2"},
	{261.63, "3"},
	{277.18, "4"},
	{293.66, "5"},
	{311.13, "6"},
	{329.63, "7"},
	{349.23, "8"},
	{369.99, "9"},
	{392.00, "10"},
	{415.30, "11"},
	{440.00, "12"},
	{466.16, "13"},
	{493.88, "14"},
	{523.25, "15"},
	{554.37, "C# (5)"},
	{587.33, "D (5)"},
	{622.25, "18"},
	{659.25, "E (5)"},
	{698.46, "F (5)"},
	{739.99, "21"},
	{783.99, "G (5)"},
	{830.61, "23"},
	{880.00, "A (5)"},
	{932.33, "25"},
	{987.77, "B/H (5)"},
	{1046.50, "C (6)"},
	{1108.73, "28"},
	{1174.66, "D (6)"},
	{1244.51, "30"},
	{1318.51, "31"},
	{1396.91, "32"},
	{1479.98, "33"},
	{1567.98, "34"},
	{1661.22, "35"},
	{1760.00, "36"},
	{1864.66, "37"},
	{1975.53, "38"},
	{2093.00, "39"},
}

func guessNote(frequency float64) note {
	min := 0
	max := len(notes) - 1
	for {
		if frequency <= notes[min].Frequency {
			return notes[min]
		}
		if frequency >= notes[max].Frequency {
			return notes[max]
		}
		if max-min <= 1 {
			toMax := notes[max].Frequency - frequency
			toMin := frequency - notes[min].Frequency
			if toMax < toMin {
				return notes[max]
			}
			return notes[min]
		}
		middle := (min + max) / 2
		if frequency <= notes[middle].Frequency {
			max = middle
		} else {
			min = middle
		}
	}
}
