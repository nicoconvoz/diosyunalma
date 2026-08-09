// Package spectral measures the periodic content of a sampled series.
//
// It exists for one experiment: the explicit formula says the deviation of the
// prime count from its smooth part is a sum of cosines in ln x, one per zeta
// zero. Listening for those cosines means projecting a series onto candidate
// frequencies — which is all a periodogram is.
package spectral

import "math"

// Periodogram returns the spectral power of y, sampled at points u, at each
// candidate frequency.
//
// The mean of y is removed first, so a constant series carries no power. The
// power at frequency f is |Σ (y_j − ȳ)·e^(−i·f·u_j)|² / n: large when the
// series oscillates at f across the sampled range, near zero otherwise. The
// sample points need not be evenly spaced.
//
// A nil result signals unusable input: an empty series or mismatched lengths.
func Periodogram(u, y, freqs []float64) []float64 {
	if len(u) == 0 || len(u) != len(y) {
		return nil
	}

	mean := 0.0
	for _, v := range y {
		mean += v
	}
	mean /= float64(len(y))

	out := make([]float64, len(freqs))
	for i, f := range freqs {
		var re, im float64
		for j, uj := range u {
			c := y[j] - mean
			s, cn := math.Sincos(f * uj)
			re += c * cn
			im += c * s
		}
		out[i] = (re*re + im*im) / float64(len(u))
	}

	return out
}
