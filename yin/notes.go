package yin

type note struct {
	Frequency float64
	Name      string
}

var notes = []note{
	{-1.0, "pause"},
	{587.33, "D (5)"},
	{554.37, "C# (5)"},
	{659.25, "E (5)"},
	{698.46, "F (5)"},
	{783.99, "G (5)"},
	{880.00, "A (5)"},
	{880.00, "A (5)"},
	{987.77, "B/H (5)"},
	{1046.50, "C (6)"},
	{1174.66, "D (6)"},
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
