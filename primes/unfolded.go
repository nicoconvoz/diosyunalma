package primes

import "math"

// Unfolded returns the gaps of seq rescaled by the local mean gap ln(p), so
// the average spacing is 1 and fluctuations become comparable across the
// whole range.
//
// This is the normalisation Montgomery applied to the zeta zeros before the
// GUE connection appeared, applied here to the primes themselves. One caveat
// travels with it: prime gaps are even integers, so the small-spacing deficit
// against a continuous law partly reflects granularity, not repulsion.
func Unfolded(seq []int) []float64 {
	if len(seq) < 2 {
		return []float64{}
	}

	out := make([]float64, len(seq)-1)
	for i := range out {
		out[i] = float64(seq[i+1]-seq[i]) / math.Log(float64(seq[i]))
	}

	return out
}
