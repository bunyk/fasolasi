package ear

import (
	"log"

	"github.com/MarkKremer/microphone"
	"github.com/faiface/beep"

	"github.com/bunyk/fasolasi/src/yin"
)

type Ear struct {
	micStream     beep.Streamer
	MicBuffer     [][2]float64
	pitchDetector yin.Yin
	Pitch         float64
}

func (e *Ear) listen() {
	go func() {
		for {
			e.micStream.Stream(e.MicBuffer)
			e.Pitch = e.pitchDetector.GetPitch2(e.MicBuffer)
			e.pitchDetector.Clean()
		}
	}()
}

func New(sampleRate, bufSize int) *Ear {
	microphone.Init() // without this you will get "PortAudio not initialized" error later

	// Create microphone stream
	micStream, _, err := microphone.OpenDefaultStream(beep.SampleRate(sampleRate), 1)
	if err != nil {
		log.Fatal(err)
	}
	var e = &Ear{
		MicBuffer:     make([][2]float64, bufSize),
		micStream:     micStream,
		pitchDetector: yin.NewYin(float64(sampleRate), bufSize, 0.05),
	}
	micStream.Start() // Start recording
	e.listen()
	return e
}
