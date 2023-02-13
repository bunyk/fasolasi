package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

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

	song, err := notes.ReadSong("songs/morrowind.txt")
	if err != nil {
		log.Fatal(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(0, win.Bounds().H()/4-10), basicAtlas)
	basicTxt.Color = colornames.Black

	game := NewGame(song)

	for !win.Closed() {
		// Input
		if win.Pressed(pixelgl.KeyA) {
			e.Pitch = 523.25
		}
		if win.Pressed(pixelgl.KeyS) {
			e.Pitch = 587.33
		}
		if win.Pressed(pixelgl.KeyD) {
			e.Pitch = 659.25
		}
		if win.Pressed(pixelgl.KeyF) {
			e.Pitch = 698.0
		}
		note, _ := notes.GuessNote(e.Pitch)

		// Processing
		game.Update(note)

		// Rendering
		win.Clear(colornames.Antiquewhite)
		soundVisualization(win, colornames.Blue, e.MicBuffer)
		drawNoteLines(win)
		drawNotes(win, game.Song, game.Played, game.Duration)

		basicTxt.Clear()
		fmt.Fprintf(basicTxt, "%d", int(game.Score*100))
		basicTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 4))
		/* Draw note played
		if name != "p" {
			imd.Color = colornames.Black
			imd.Push(pixel.V(WIDTH/2, ybase+float64(line)*noteRadius*2))
			imd.Circle(noteRadius, 3)
		}
		*/
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

type Game struct {
	Song       []notes.SongNote
	Played     []playedNote
	Score      float64
	Start      time.Time
	Duration   float64 // in seconds
	SongCursor int     // number of passsed notes in song
}

type playedNote struct {
	notes.SongNote
	Correct bool
}

func NewGame(song []notes.SongNote) *Game {
	return &Game{
		Played: make([]playedNote, 0, 100),
		Song:   song,
		Start:  time.Now(),
	}
}

func (g *Game) currentNote() notes.Pitch {
	for g.SongCursor < len(g.Song) && g.Song[g.SongCursor].End() < g.Duration { // skip played notes
		g.SongCursor += 1
	}
	if g.SongCursor >= len(g.Song) {
		return notes.Pause
	}
	if g.Song[g.SongCursor].Time < g.Duration { // already should be playing
		return g.Song[g.SongCursor].Pitch
	}
	return notes.Pause
}

func (g *Game) Update(note notes.Pitch) {
	oldDuration := g.Duration
	g.Duration = time.Since(g.Start).Seconds()
	playingCorrectly := note == g.currentNote()
	if playingCorrectly {
		g.Score += g.Duration - oldDuration
	}
	if note.Name != "p" {
		if len(g.Played) == 0 || g.Played[len(g.Played)-1].End() > 0 { // no note currently playing
			g.Played = append(g.Played, playedNote{
				SongNote: notes.SongNote{ // create new note
					Time:     g.Duration,
					Pitch:    note,
					Duration: -1.0,
				},
				Correct: playingCorrectly,
			})
		} else if g.Played[len(g.Played)-1].Pitch != note || g.Played[len(g.Played)-1].Correct != playingCorrectly { // note changed
			g.Played[len(g.Played)-1].Duration = g.Duration - g.Played[len(g.Played)-1].Time // end current one
			g.Played = append(g.Played, playedNote{
				SongNote: notes.SongNote{ // create new note
					Time:     g.Duration,
					Pitch:    note,
					Duration: -1.0,
				},
				Correct: playingCorrectly,
			})
		}
	} else { // no note
		if len(g.Played) > 0 && g.Played[len(g.Played)-1].End() < 0 { // there is a note still playing
			g.Played[len(g.Played)-1].Duration = g.Duration - g.Played[len(g.Played)-1].Time // end it
		}
	}

	if len(g.Played) > 2 && g.Played[0].End() < g.Duration-timeLinePosition/noteSPS { // note not visible
		g.Played = g.Played[1:] // remove
	}
}

const noteRadius = 40
const noteSPS = 0.15 // screens per second
const timeLinePosition = 0.3

func time2X(width, time, currentTime float64) float64 {
	return width * (timeLinePosition - (currentTime-time)*noteSPS)
}

func drawNotes(win *pixelgl.Window, song []notes.SongNote, played []playedNote, time float64) {
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.SharpEndShape
	width := win.Bounds().W()
	ybase := float64(win.Bounds().H()/2 - noteRadius*4)

	for _, note := range song {
		drawNote(imd, time, width, ybase, false, note)
	}
	for _, note := range played {
		drawNote(imd, time, width, ybase, note.Correct, note.SongNote)
	}
	imd.Draw(win)
}

var rainbow = []color.Color{
	colornames.Red,
	colornames.Orange,
	colornames.Yellow,
	colornames.Green,
	colornames.Blue,
	colornames.Violet,
}

func drawNote(imd *imdraw.IMDraw, time, width, ybase float64, colorful bool, note notes.SongNote) {
	if note.Pitch.Name == "p" {
		return
	}
	end := note.End()
	if end < 0.0 {
		end = time
	}
	endX := time2X(width, end, time)
	if endX < 0 { // invisible already
		return
	}
	startX := time2X(width, note.Time, time)
	if startX > width { // still invisible
		return
	}
	ycenter := ybase + note.Pitch.Bottom*noteRadius*2

	if note.Pitch.HasAdditionalLine() {
		imd.Color = colornames.Black
		imd.Push(
			pixel.V(startX-noteRadius*2, ycenter),
			pixel.V(endX+noteRadius*2, ycenter),
		)
		imd.Line(1)
	}
	border := 0.0
	if colorful {
		imd.Color = rainbow[(10+int(note.Pitch.Bottom*4))%len(rainbow)]
	} else {
		if note.Pitch.Height >= 0.5 { // white key
			imd.Color = colornames.White
			imd.Push(
				pixel.V(startX+1, ycenter-note.Pitch.Height*noteRadius+1),
				pixel.V(endX-1, ycenter+note.Pitch.Height*noteRadius-1),
			)
			imd.Rectangle(0)
			border = 2.0
		}
		imd.Color = colornames.Black
	}

	imd.Push(
		pixel.V(startX+1, ycenter-note.Pitch.Height*noteRadius+1),
		pixel.V(endX-1, ycenter+note.Pitch.Height*noteRadius-1),
	)
	imd.Rectangle(border)
}

func drawNoteLines(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	width := win.Bounds().W()
	height := win.Bounds().H()

	ybase := float64(height/2 - noteRadius*4)
	for i := 0; i < 5; i++ {
		imd.Push(
			pixel.V(0, ybase+float64(i)*noteRadius*2),
			pixel.V(width, ybase+float64(i)*noteRadius*2),
		)
		imd.Line(1)
	}
	imd.Push(
		pixel.V(width*timeLinePosition, 0),
		pixel.V(width*timeLinePosition, height),
	)
	imd.Line(3)

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
