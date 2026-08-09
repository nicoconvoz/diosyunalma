package riemann

import "math"

// ZeroCount returns the Riemann–von Mangoldt estimate of how many nontrivial
// zeta zeros have height between 0 and t:
//
//	N(t) = (t/2π)·ln(t/2πe) + 7/8
//
// The exact count adds a fluctuation term S(t) that stays below 1 in this
// range, so the smooth formula is accurate to better than one zero.
//
// The same expression is the semiclassical state count of the Berry–Keating
// Hamiltonian H = xp — which is why it matters here. If a Hilbert–Pólya
// operator exists, this is its density of states, and a measured zero census
// must follow it.
func ZeroCount(t float64) float64 {
	if t <= 0 {
		return 0
	}
	return t/(2*math.Pi)*math.Log(t/(2*math.Pi*math.E)) + 7.0/8.0
}
