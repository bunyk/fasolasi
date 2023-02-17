package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/ear"
	"github.com/bunyk/fasolasi/src/game"
	"github.com/bunyk/fasolasi/src/notes"
)

func run() {
	song, err := notes.ReadSong(config.SongFilename)
	if err != nil {
		log.Fatal(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:     "FaSoLaSi",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// FPS tracking
	frames := 0
	second := time.Tick(time.Second)

	// Audio input stuff
	e := ear.New(config.MicrophoneSampleRate, config.MicrophoneBufferLength)

	gameSession := game.NewSession(song, config.GameMode)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Input
		kp := KeyboardPitch(win)
		if kp > 0.0 {
			e.Pitch = kp
		}
		note, _ := notes.GuessNote(e.Pitch)

		// Processing
		gameSession.Update(dt, note)

		// Rendering
		win.Clear(colornames.Antiquewhite)
		soundVisualization(win, colornames.Blue, e.MicBuffer)
		gameSession.Render(win)

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

func KeyboardPitch(win *pixelgl.Window) float64 {
	if win.Pressed(pixelgl.KeyA) {
		return 523.25
	}
	if win.Pressed(pixelgl.KeyS) {
		return 587.33
	}
	if win.Pressed(pixelgl.KeyD) {
		return 659.25
	}
	if win.Pressed(pixelgl.KeyF) {
		return 698.0
	}
	return -1.0
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
