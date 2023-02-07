package main

import (
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/MarkKremer/microphone"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
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

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		buf := visualizer.Buffer()

		lineGraph(win, colornames.Blue, buf)
		win.Update()
	}
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

type Visualizer struct {
	mu     sync.RWMutex
	buffer []float64
	cursor int
	S      beep.Streamer
}

func (v *Visualizer) Buffer() []float64 {
	v.mu.RLock()
	res := make([]float64, len(v.buffer))
	for i := range v.buffer {
		res[i] = v.buffer[(v.cursor+i)%len(v.buffer)]
	}
	v.mu.RUnlock()
	return res
}

func (v *Visualizer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = v.S.Stream(samples)
	v.mu.Lock()
	for _, s := range samples {
		v.buffer[v.cursor] = s[0]
		v.cursor = (v.cursor + 1) % len(v.buffer)
	}
	v.mu.Unlock()
	return
}

func (v Visualizer) Err() error {
	return v.S.Err()
}

var visualizer = &Visualizer{
	buffer: make([]float64, 11025),
}

func main() {
	// First, some configuration
	sampleRate := beep.SampleRate(44100) // Choose sample rate
	numChannels := 1                     // 1 - mono, 2 - stereo

	microphone.Init() // without this you will get "PortAudio not initialized" error later

	// Create microphone stream
	micStream, format, err := microphone.OpenDefaultStream(sampleRate, numChannels)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", format) // Sample rate and number of channels, in case you'll need it

	speaker.Init(sampleRate, sampleRate.N(time.Second)) // Initialize speaker, with a buffer for 0.1 sec

	micStream.Start() // Start recording
	visualizer.S = micStream
	// speaker.Play(visualizer) // Start playing what you record

	go func() {
		samples := make([][2]float64, 512)
		for {
			visualizer.Stream(samples)
		}
	}()

	pixelgl.Run(run)
}
