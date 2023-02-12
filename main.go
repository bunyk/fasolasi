package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bunyk/fasolasi/src/ear"
	"github.com/bunyk/fasolasi/src/notes"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "FaSoLaSi",
		Bounds:    pixel.R(0, 0, 1024, 768),
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
	const bufLen = 11025 / 2
	const sampleRate = 188200
	e := ear.New(sampleRate, bufLen)

	start := time.Now()
	played := make([]playedNote, 0, 100)

	song, err := notes.ReadSong("songs/range.txt")
	for _, n := range song {
		fmt.Printf("%#v\n", n)
	}

	for !win.Closed() {
		t := time.Since(start).Seconds()
		note, nn := notes.GuessNote(e.Pitch)
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

		win.Clear(colornames.Antiquewhite)

		soundVisualization(win, colornames.Blue, e.MicBuffer)
		drawLines(win)
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
const notePPS = 150 // pixels per second

func time2X(width, time, currentTime float64) float64 {
	return width/2 - (currentTime-time)*notePPS
}

var rainbow = []color.Color{
	colornames.Red,
	colornames.Orange,
	colornames.Yellow,
	colornames.Green,
	colornames.Blue,
	colornames.Violet,
}

func drawSongNotes(win *pixelgl.Window, song []notes.SongNote, time float64) {
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.SharpEndShape
	imd.Color = colornames.Black
	width := win.Bounds().W()
	base := float64(win.Bounds().H()/2 - noteRadius*4)

	for _, note := range song {
		if note.Pitch.Name == "p" {
			continue
		}
		if note.End() > 0 && time2X(width, note.End(), time) < 0 { // invisible already
			continue
		}
		startX := time2X(width, note.Time, time)
		endX := time2X(width, note.End(), time)
		ycenter := base + note.Pitch.Bottom*noteRadius*2
		if note.Pitch.HasAdditionalLine() {
			imd.Push(
				pixel.V(startX-noteRadius*2, ycenter),
				pixel.V(endX+noteRadius*2, ycenter),
			)
			imd.Line(1)
		}
		imd.Push(
			pixel.V(startX+1, ycenter-note.Pitch.Height*noteRadius+1),
			pixel.V(endX-1, ycenter+note.Pitch.Height*noteRadius-1),
		)
		border := 2.0
		if note.Pitch.Height < 0.5 { // black key
			border = 0.0
		}
		imd.Rectangle(border)
	}
	imd.Draw(win)
}

func drawLines(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	width := win.Bounds().W()
	height := win.Bounds().H()

	base := float64(height/2 - noteRadius*4)
	for i := 0; i < 5; i++ {
		imd.Push(
			pixel.V(0, base+float64(i)*noteRadius*2),
			pixel.V(width, base+float64(i)*noteRadius*2),
		)
		imd.Line(1)
	}

	imd.Draw(win)
}

func drawNotes(win *pixelgl.Window, line float64, name string, played []playedNote, time float64) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	y := float64(win.Bounds().H()/2 - noteRadius*4)
	width := win.Bounds().W()

	for _, note := range played {
		if note.end > 0 && time2X(width, note.end, time) < 0 { // invisible already
			continue
		}
		end := note.end
		if end == 0 {
			end = time
		}
		imd.Color = rainbow[(10+int(note.line*4))%len(rainbow)]
		imd.Push(
			pixel.V(time2X(width, note.start, time), y+(float64(note.line)-0.25)*noteRadius*2),
			pixel.V(time2X(width, end, time), y+(float64(note.line)+0.25)*noteRadius*2),
		)
		imd.Rectangle(0)
	}
	/* Draw note played
	if name != "p" {
		imd.Color = colornames.Black
		imd.Push(pixel.V(WIDTH/2, y+float64(line)*noteRadius*2))
		imd.Circle(noteRadius, 3)
	}
	*/
	imd.Draw(win)
}

func soundVisualization(win *pixelgl.Window, col color.Color, data [][2]float64) {
	imd := imdraw.New(nil)
	imd.Color = col
	width := win.Bounds().W()
	height := win.Bounds().H()
	every := 1.0
	if width > 0 {
		every = float64(len(data)) / width
	}
	for i := 0.0; i < width; i += 1.0 {
		imd.Push(pixel.V(float64(i), height*(0.1+data[int(i*every)][0]*0.2)))
	}
	imd.Line(1)
	imd.Draw(win)
}

func main() {
	pixelgl.Run(run)
}
