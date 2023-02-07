package ear

import (
	"fmt"
	"log"
	"sync"

	"github.com/MarkKremer/microphone"
	"github.com/faiface/beep"
)

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

func NewVisualizer(bufSize, sampleRate int) *Visualizer {
	microphone.Init() // without this you will get "PortAudio not initialized" error later

	// Create microphone stream
	numChannels := 1 // 1 - mono, 2 - stereo
	micStream, format, err := microphone.OpenDefaultStream(beep.SampleRate(sampleRate), numChannels)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", format) // Sample rate and number of channels, in case you'll need it
	var visualizer = &Visualizer{
		buffer: make([]float64, bufSize),
		S:      micStream,
	}
	micStream.Start() // Start recording
	go func() {
		samples := make([][2]float64, 512)
		for {
			visualizer.Stream(samples)
		}
	}()
	return visualizer
}
