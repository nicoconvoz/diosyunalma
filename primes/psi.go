package primes

import (
	"math"
	"sort"
)

// PsiAt returns Chebyshev's psi function evaluated at each of xs, which must be
// ascending. psi(x) sums ln p over every prime power p^k up to x.
//
// psi is the natural observable for listening to the zeta zeros: the explicit
// formula gives psi(x) = x − Σ_ρ x^ρ/ρ − …, so (psi(x) − x)/√x is, up to small
// corrections, a superposition of cosines in ln x — one per zero. The prime
// powers matter: dropping them distorts exactly the low-x range where the
// sampling grid is densest.
func PsiAt(limit int, xs []int) []float64 {
	if len(xs) == 0 {
		return []float64{}
	}

	type event struct {
		at     int
		weight float64
	}

	// Squares and higher powers are few — about √limit of them — so collecting
	// and sorting them separately keeps the prime list untouched.
	powers := []event{}
	ps := Sieve(limit)
	for _, p := range ps {
		if p*p > limit {
			break
		}
		w := math.Log(float64(p))
		for pk := p * p; pk <= limit; pk *= p {
			powers = append(powers, event{pk, w})
		}
	}
	sort.Slice(powers, func(i, j int) bool { return powers[i].at < powers[j].at })

	out := make([]float64, len(xs))
	sum := 0.0
	pi, wi := 0, 0
	for i, x := range xs {
		for pi < len(ps) && ps[pi] <= x {
			sum += math.Log(float64(ps[pi]))
			pi++
		}
		for wi < len(powers) && powers[wi].at <= x {
			sum += powers[wi].weight
			wi++
		}
		out[i] = sum
	}

	return out
}
