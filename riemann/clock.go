package riemann

import "math"

// PrimeClock evaluates the truncated explicit formula for the prime density at
// x, using the given zero heights as clock hands:
//
//	1 − (2/√x) · Σⱼ cos(γⱼ · ln x)
//
// Each zero is a hand rotating at angular speed γⱼ as u = ln x advances. The
// full formula, with every zero, is a sum of delta spikes at the prime powers:
// the hands align exactly there. A truncated set of hands gives a smoothed
// version whose local maxima still mark the prime powers — which turns a
// measured zero list into a device that points back at the primes.
//
// This is the inverse of the zeta hunt. That measurement ran primes → zeros;
// the clock runs zeros → primes, and the two directions together are the
// duality the explicit formula asserts.
func PrimeClock(gammas []float64, x float64) float64 {
	if x <= 1 {
		return 0
	}

	u := math.Log(x)
	sum := 0.0
	for _, g := range gammas {
		sum += math.Cos(g * u)
	}

	return 1 - 2*sum/math.Sqrt(x)
}
