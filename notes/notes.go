package notes

type note struct {
	Frequency float64
	Name      string
	Line      int // even - on line, odd - between lines
}

var notes = []note{
	{-1.0, "pause", 0},
	{523.25, "C (5)", -2},
	{554.37, "C# (5)", -2},
	{587.33, "D (5)", -1},
	{622.25, "D# (5)", -1},
	{659.25, "E (5)", 0},
	{698.46, "F (5)", 1},
	{739.99, "F#", 1},
	{783.99, "G (5)", 2},
	{830.61, "G# (5)", 2},
	{880.00, "A (5)", 3},
	{932.33, "Bb/Hb", 4},
	{987.77, "B/H (5)", 4},
	{1046.50, "C (6)", 5},
	{1108.73, "C# (6)", 5},
	{1174.66, "D (6)", 6},
	{1244.51, "D# (6)", 6},
	{1318.51, "E (6)", 7},
	{1396.91, "F (6)", 8},
	{1479.98, "33", 10}, // TODO
	{1567.98, "34", 10},
	{1661.22, "35", 10},
	{1760.00, "36", 10},
	{1864.66, "37", 10},
	{1975.53, "38", 10},
	{2093.00, "39", 10},
}

func GuessNote(frequency float64) (note, int) {
	min := 0
	max := len(notes) - 1
	for {
		if frequency <= notes[min].Frequency {
			return notes[min], min
		}
		if frequency >= notes[max].Frequency {
			return notes[max], max
		}
		if max-min <= 1 {
			toMax := notes[max].Frequency - frequency
			toMin := frequency - notes[min].Frequency
			if toMax < toMin {
				return notes[max], max
			}
			return notes[min], min
		}
		middle := (min + max) / 2
		if frequency <= notes[middle].Frequency {
			max = middle
		} else {
			min = middle
		}
	}
}
