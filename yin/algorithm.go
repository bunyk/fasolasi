// Pitch Detection Algorithm
//
// Given an audio sample buffer detects the main pitch of the
//	audio, and the probability of it being correct.
//	With this information you can detect the note that is played.
//
// The frequency detection algorithm is based on YIN algorithm more specifically
// a port to the Go ( GoLang ) programming language of the of the C implementation
// on https://github.com/ashokfernandez/Yin-Pitch-Tracking/blob/master/Yin.c
//
// Author:  Joao Nuno Carvalho
// Email:   joaonunocarv@gmail.com
// Date:    2017.12.9
// License: MIT OpenSource License

package yin

type Yin struct {
	samplingRate   float64
	bufferSize     int       // Size of the buffer to process.
	halfBufferSize int       // Half of buffer size.
	yinBuffer      []float64 // Buffer that stores the results of the intermediate processing steps of the algorithm
	probability    float64   // Probability that the pitch found is correct as a decimal (i.e 0.85 is 85%)
	threshold      float64   // Allowed uncertainty in the result as a decimal (i.e 0.15 is 15%)
}

// threshold  - Allowed uncertainty (e.g 0.05 will return a pitch with ~95% probability)
func NewYin(samplingRate float64, bufferSize int, threshold float64) Yin {
	return Yin{
		samplingRate:   samplingRate,
		bufferSize:     bufferSize,
		halfBufferSize: bufferSize / 2,
		yinBuffer:      make([]float64, bufferSize/2),
		threshold:      threshold,
	}
}

func (y *Yin) Clean() {
	y.probability = 0.0

	// Allocate the autocorellation buffer and initialise it to zero.
	for i := range y.yinBuffer {
		y.yinBuffer[i] = 0
	}
}

// Runs the Yin pitch detection algortihm
//
//	buffer       - Buffer of samples to analyse
//
// return pitchInHertz - Fundamental frequency of the signal in Hz. Returns -1 if pitch can't be found
func (y *Yin) GetPitch(buffer []float64) (pitchInHertz float64) {
	//tauEstimate int      := -1
	pitchInHertz = -1

	// Step 1: Calculates the squared difference of the signal with a shifted version of itself.
	y.yinDifference(buffer)

	// Step 2: Calculate the cumulative mean on the normalised difference calculated in step 1.
	y.yinCumulativeMeanNormalizedDifference()

	// Step 3: Search through the normalised cumulative mean array and find values that are over the threshold.
	tauEstimate := y.yinAbsoluteThreshold()

	// Step 5: Interpolate the shift value (tau) to improve the pitch estimate.
	if tauEstimate != -1 {
		pitchInHertz = y.samplingRate / y.yinParabolicInterpolation(tauEstimate)
	}

	return pitchInHertz
}

// Certainty of the pitch found
// return ptobability - Returns the certainty of the note found as a decimal (i.e 0.3 is 30%)
func (y *Yin) GetProbability() (probability float64) {
	return y.probability
}

// Step 1: Calculates the squared difference of the signal with a shifted version of itself.
// @param buffer Buffer of samples to process.
//
// This is the Yin algorithms tweak on autocorellation. Read http://audition.ens.fr/adc/pdf/2002_JASA_YIN.pdf
// for more details on what is in here and why it's done this way.
func (y *Yin) yinDifference(buffer []float64) {
	// Calculate the difference for difference shift values (tau) for the half of the samples.
	for tau := 0; tau < y.halfBufferSize; tau++ {

		// Take the difference of the signal with a shifted version of itself, then square it.
		// (This is the Yin algorithm's tweak on autocorellation)
		for i := 0; i < y.halfBufferSize; i++ {
			delta := float64(buffer[i]) - float64(buffer[i+tau])
			y.yinBuffer[tau] += delta * delta
		}
	}
}

// Step 2: Calculate the cumulative mean on the normalised difference calculated in step 1
//
// This goes through the Yin autocorellation values and finds out roughly where shift is which
// produced the smallest difference
func (y *Yin) yinCumulativeMeanNormalizedDifference() {
	runningSum := 0.0
	y.yinBuffer[0] = 1

	// Sum all the values in the autocorellation buffer and nomalise the result, replacing
	// the value in the autocorellation buffer with a cumulative mean of the normalised difference.
	for tau := 1; tau < y.halfBufferSize; tau++ {
		runningSum += y.yinBuffer[tau]
		y.yinBuffer[tau] *= float64(tau) / runningSum
	}
}

// Step 3: Search through the normalised cumulative mean array and find values that are over the threshold
// return Shift (tau) which caused the best approximate autocorellation. -1 if no suitable value is found
// over the threshold.
func (y *Yin) yinAbsoluteThreshold() int {

	var tau int

	// Search through the array of cumulative mean values, and look for ones that are over the threshold
	// The first two positions in yinBuffer are always so start at the third (index 2)
	for tau = 2; tau < y.halfBufferSize; tau++ {
		if y.yinBuffer[tau] < y.threshold {
			for (tau+1 < y.halfBufferSize) && (y.yinBuffer[tau+1] < y.yinBuffer[tau]) {
				tau++
			}

			/* found tau, exit loop and return
			 * store the probability
			 * From the YIN paper: The yin->threshold determines the list of
			 * candidates admitted to the set, and can be interpreted as the
			 * proportion of aperiodic power tolerated
			 * within a periodic signal.
			 *
			 * Since we want the periodicity and and not aperiodicity:
			 * periodicity = 1 - aperiodicity */
			y.probability = 1 - y.yinBuffer[tau]
			break
		}
	}

	// if no pitch found, tau => -1
	if tau == y.halfBufferSize || y.yinBuffer[tau] >= y.threshold {
		tau = -1
		y.probability = 0
	}

	return tau
}

// Step 5: Interpolate the shift value (tau) to improve the pitch estimate.
// tauEstimate [description]
// Return
// The 'best' shift value for autocorellation is most likely not an interger shift of the signal.
// As we only autocorellated using integer shifts we should check that there isn't a better fractional
// shift value.
func (y *Yin) yinParabolicInterpolation(tauEstimate int) float64 {

	var betterTau float64
	var x0 int
	var x2 int

	// Calculate the first polynomial coeffcient based on the current estimate of tau.
	if tauEstimate < 1 {
		x0 = tauEstimate
	} else {
		x0 = tauEstimate - 1
	}

	// Calculate the second polynomial coeffcient based on the current estimate of tau.
	if tauEstimate+1 < y.halfBufferSize {
		x2 = tauEstimate + 1
	} else {
		x2 = tauEstimate
	}

	// Algorithm to parabolically interpolate the shift value tau to find a better estimate.
	if x0 == tauEstimate {
		if y.yinBuffer[tauEstimate] <= y.yinBuffer[x2] {
			betterTau = float64(tauEstimate)
		} else {
			betterTau = float64(x2)
		}
	} else if x2 == tauEstimate {
		if y.yinBuffer[tauEstimate] <= y.yinBuffer[x0] {
			betterTau = float64(tauEstimate)
		} else {
			betterTau = float64(x0)
		}
	} else {
		var s0, s1, s2 float64
		s0 = y.yinBuffer[x0]
		s1 = y.yinBuffer[tauEstimate]
		s2 = y.yinBuffer[x2]
		// fixed AUBIO implementation, thanks to Karl Helgason:
		// (2.0f * s1 - s2 - s0) was incorrectly multiplied with -1
		betterTau = float64(tauEstimate) + (s2-s0)/(2*(2*s1-s2-s0))
	}

	return betterTau
}
