package yin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unixpickle/wav"
)

func TestPitchDetection(t *testing.T) {
	s, err := wav.ReadSoundFile("test.wav")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Sample rate:  %d\n", s.SampleRate()) // 44100

	const testBuffLen int = 11025
	yin := NewYin(float64(s.SampleRate()), testBuffLen, 0.05)

	// frequencies and their probabilities in test.wav file
	var expected = [][2]float64{
		{-1.000000, 0.000000},
		{-1.000000, 0.000000},
		{-1.000000, 0.000000},
		{-1.000000, 0.000000},
		{785.740730, 0.971132},
		{798.119118, 0.957630},
		{883.167428, 0.982349},
		{983.579367, 0.997792},
		{982.246827, 0.995798},
	}

	samples := make([]float64, len(s.Samples()))
	for i, s := range s.Samples() { // it's a slice of wav.Sample, which is just alias for float64. Why?!?!
		samples[i] = float64(s)
	}
	for i, exp := range expected {
		// Process the buffer with the algorithm of YIN for frequency detection.
		frequency := yin.GetPitch(samples[i*testBuffLen : (i+1)*testBuffLen])
		probability := yin.GetProbability()
		yin.Clean()
		assert.InDelta(t, exp[0], frequency, 1e-6)
		assert.InDelta(t, exp[1], probability, 1e-6)
	}
}
