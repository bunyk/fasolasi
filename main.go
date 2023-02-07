package main

import (
	"image/color"

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
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Audio input stuff
	const bufLen = 11025
	const sampleRate = 44100
	visualizer := ear.NewVisualizer(bufLen, 44100)
	detector := yin.NewYin(sampleRate, bufLen, 0.05)

	for !win.Closed() {

		buf := visualizer.Buffer()
		pitch := detector.GetPitch(buf)
		detector.Clean()
		_, nn := notes.GuessNote(pitch)

		win.Clear(colornames.Whitesmoke)

		lineGraph(win, colornames.Blue, buf)
		drawNote(win, nn)
		win.Update()
	}
}

func drawNote(t pixel.Target, num int) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	for i := 1; i < 50; i += 2 {
		imd.Push(
			pixel.V(0, float64(i*HEIGHT/50)),
			pixel.V(WIDTH, float64(i*HEIGHT/50)),
		)
		imd.Line(1)
	}
	if num > 0 {
		imd.Push(pixel.V(WIDTH/2, float64(num*HEIGHT/50)))
		imd.Circle(30, 5)
	}
	imd.Draw(t)
}

func lineGraph(t pixel.Target, col color.Color, data []float64) {
	imd := imdraw.New(nil)
	imd.Color = col
	for i := 0; i < len(data)/5; i++ {
		imd.Push(pixel.V(float64(i)*0.5, HEIGHT*(0.5+data[i*5])))
	}
	imd.Line(1)
	imd.Draw(t)
}

func main() {
	pixelgl.Run(run)
}
