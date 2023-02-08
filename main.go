package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/bunyk/fasolasi/ear"
	"github.com/bunyk/fasolasi/notes"
	"github.com/bunyk/fasolasi/yin"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const BUFFER_SIZE = 512
const WIDTH = 1024
const HEIGHT = 768

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "FaSoLaSi",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
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
	const sampleRate = 44100 / 2
	visualizer := ear.NewVisualizer(bufLen, sampleRate)
	detector := yin.NewYin(sampleRate, bufLen, 0.05)

	start := time.Now()
	played := make([]playedNote, 0, 100)
	for !win.Closed() {
		t := time.Since(start).Seconds()
		buf := visualizer.Buffer()
		pitch := detector.GetPitch(buf)
		detector.Clean()
		note, nn := notes.GuessNote(pitch)
		if nn > 0 {
			if len(played) == 0 || played[len(played)-1].end > 0 { // no note currently playing
				played = append(played, playedNote{ // create new note
					start: t,
					line:  note.Line,
				})
			} else if played[len(played)-1].line != note.Line { // note changed
				played[len(played)-1].end = t       // end current one
				played = append(played, playedNote{ // create new note
					start: t,
					line:  note.Line,
				})
			}
		} else { // no note
			if len(played) > 0 && played[len(played)-1].end == 0 { // there is a note
				played[len(played)-1].end = t // end it
			}
		}

		win.Clear(colornames.Whitesmoke)

		lineGraph(win, colornames.Blue, buf)
		drawNotes(win, note.Line, note.Name, played, t)
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
	line  int
}

const noteRadius = 20
const notePPS = 100 // pixels per second

func time2X(time, currentTime float64) float64 {
	return WIDTH/2 - (currentTime-time)*notePPS
}

func drawNotes(t pixel.Target, line int, name string, played []playedNote, time float64) {
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
	imd.Color = colornames.Green
	for _, note := range played {
		if note.end > 0 && time2X(note.end, time) < 0 { // invisible already
			continue
		}
		end := note.end
		if end == 0 {
			end = time
		}
		imd.Push(
			pixel.V(time2X(note.start, time), y+(float64(note.line)-0.5)*noteRadius),
			pixel.V(time2X(end, time), y+(float64(note.line)+0.5)*noteRadius),
		)
		imd.Rectangle(0)
	}
	if name != "pause" {
		imd.Color = colornames.Black
		imd.Push(pixel.V(WIDTH/2, y+float64(line)*noteRadius))
		imd.Circle(noteRadius, 0)
	}
	imd.Draw(t)
}

func lineGraph(t pixel.Target, col color.Color, data []float64) {
	imd := imdraw.New(nil)
	imd.Color = col
	for i := 0; i < WIDTH; i++ {
		imd.Push(pixel.V(float64(i), HEIGHT*(0.1+data[i*len(data)/WIDTH])))
	}
	imd.Line(1)
	imd.Draw(t)
}

func main() {
	pixelgl.Run(run)
}
