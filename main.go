package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/bunyk/fasolasi/ear"
	"github.com/bunyk/fasolasi/notes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const WIDTH = 1024
const HEIGHT = 768

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "FaSoLaSi",
		Bounds:    pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// FPS tracking
	frames := 0
	second := time.Tick(time.Second)

	// Audio input stuff
	const bufLen = 11025 / 4
	const sampleRate = 88200
	e := ear.New(sampleRate, bufLen)

	start := time.Now()
	played := make([]playedNote, 0, 100)

	song, err := notes.ReadSong("songs/range.txt")
	fmt.Println(song, err)

	for !win.Closed() {
		t := time.Since(start).Seconds()
		pitch := e.Pitch
		note, nn := notes.GuessNote(pitch)
		if nn > 0 {
			if len(played) == 0 || played[len(played)-1].end > 0 { // no note currently playing
				played = append(played, playedNote{ // create new note
					start: t,
					line:  note.Bottom,
				})
			} else if played[len(played)-1].line != note.Bottom { // note changed
				played[len(played)-1].end = t       // end current one
				played = append(played, playedNote{ // create new note
					start: t,
					line:  note.Bottom,
				})
			}
		} else { // no note
			if len(played) > 0 && played[len(played)-1].end == 0 { // there is a note
				played[len(played)-1].end = t // end it
			}
		}

		win.Clear(colornames.Whitesmoke)

		lineGraph(win, colornames.Blue, e.MicBuffer)
		drawNotes(win, note.Bottom, note.Name, played, t)
		drawSongNotes(win, song, t)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("FaSoLaSi | FPS: %d", frames))
			frames = 0
		default:
		}
	}
}

type playedNote struct {
	start float64
	end   float64
	line  float64
}

const noteRadius = 40
const notePPS = 100 // pixels per second

func time2X(time, currentTime float64) float64 {
	return WIDTH/2 - (currentTime-time)*notePPS
}

var rainbow = []color.Color{
	colornames.Red,
	colornames.Orange,
	colornames.Yellow,
	colornames.Green,
	colornames.Blue,
	colornames.Violet,
}

func drawSongNotes(t pixel.Target, song []notes.SongNote, time float64) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	base := float64(HEIGHT/2 - noteRadius*4)

	for _, note := range song {
		if note.End() > 0 && time2X(note.End(), time) < 0 { // invisible already
			continue
		}
		ycenter := base + note.Pitch.Bottom*noteRadius*2
		imd.Push(
			pixel.V(time2X(note.Time, time), ycenter-note.Pitch.Height*noteRadius),
			pixel.V(time2X(note.End(), time), ycenter+note.Pitch.Height*noteRadius),
		)
		border := 3.0
		if note.Pitch.Height < 0.5 { // black key
			border = 0.0
		}
		imd.Rectangle(border)
	}
	imd.Draw(t)
}

func drawNotes(t pixel.Target, line float64, name string, played []playedNote, time float64) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	y := float64(HEIGHT/2 - noteRadius*4)
	for i := 0; i <= 8; i += 2 {
		imd.Push(
			pixel.V(0, y+float64(i)*noteRadius),
			pixel.V(WIDTH, y+float64(i)*noteRadius),
		)
		imd.Line(1)
	}

	for _, note := range played {
		if note.end > 0 && time2X(note.end, time) < 0 { // invisible already
			continue
		}
		end := note.end
		if end == 0 {
			end = time
		}
		imd.Color = rainbow[(10+int(note.line*4))%len(rainbow)]
		imd.Push(
			pixel.V(time2X(note.start, time), y+(float64(note.line)-0.25)*noteRadius*2),
			pixel.V(time2X(end, time), y+(float64(note.line)+0.25)*noteRadius*2),
		)
		imd.Rectangle(0)
	}
	if name != "pause" {
		imd.Color = colornames.Black
		imd.Push(pixel.V(WIDTH/2, y+float64(line)*noteRadius*2))
		imd.Circle(noteRadius, 3)
	}
	imd.Draw(t)
}

func lineGraph(t pixel.Target, col color.Color, data [][2]float64) {
	imd := imdraw.New(nil)
	imd.Color = col
	for i := 0; i < WIDTH; i++ {
		imd.Push(pixel.V(float64(i), HEIGHT*(0.1+data[i*len(data)/WIDTH][0]*0.2)))
	}
	imd.Line(1)
	imd.Draw(t)
}

func main() {
	pixelgl.Run(run)
}
