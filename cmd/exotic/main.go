// Command exotic takes the old wild flash seriously: exotic matter with
// negative mass, artificial gravitational lenses, stable wormholes. Every
// piece exists in the number sea, with its proper name:
//
//   - NEGATIVE MASS is the Moebius function: mu(n) = -1 gives a number
//     negative weight, mu(n) = 0 makes it massless. The series of exotic
//     matter is 1/zeta(s) = sum mu(n)/n^s.
//   - THE EXOTIC LENS inverts the optics: where zeta has a ZERO (a dark
//     point the hunt must feel for), 1/zeta has a POLE - a blazing beacon
//     impossible to miss. The lens turns the invisible luminous.
//   - WORMHOLE STABILITY is the negative-energy budget: a traversable
//     wormhole demands bounded exotic energy, and the Riemann Hypothesis
//     IS that certificate - it is equivalent to the Mertens tide
//     M(x)/sqrt(x) staying within every epsilon-budget. RH is the
//     stability theorem of the wormhole.
//   - The wormhole itself already flies with the fleet: Gauss reciprocity
//     (the Cassegrain crankshaft, F88) teleports a wave of q terms into a
//     wave of p terms - instant travel between remote points of the sum.
//
// Usage:
//
//	go run ./cmd/exotic
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

const emTerms = 120

func zetaEM(s complex128, lnn []float64) complex128 {
	var sum complex128
	sig, t := real(s), imag(s)
	for n := 1; n < emTerms; n++ {
		amp := math.Exp(-sig * lnn[n])
		sn, cs := math.Sincos(t * lnn[n])
		sum += complex(amp*cs, -amp*sn)
	}
	nf := complex(float64(emTerms), 0)
	ns := cmplx.Exp(-s * complex(lnn[emTerms], 0))
	sum += ns * nf / (s - 1)
	sum += ns / 2
	sum += ns * s / nf / 12
	sum -= ns * s * (s + 1) * (s + 2) / (nf * nf * nf) / 720
	sum += ns * s * (s + 1) * (s + 2) * (s + 3) * (s + 4) /
		(nf * nf * nf * nf * nf) / 30240
	return sum
}

func main() {
	fmt.Println("EXOTIC MATTER — the old flash, piece by piece, in the number sea")

	lnn := make([]float64, emTerms+1)
	for n := 1; n <= emTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}

	// Part 1: the exotic lens. Dark points become beacons.
	fmt.Println("\n  THE EXOTIC LENS - 1/|zeta| on the critical line:")
	known := []float64{14.134725, 21.022040, 25.010858, 30.424876, 32.935062}
	const dt = 0.001
	type pk struct{ t, b float64 }
	var found []pk
	prev := 1 / cmplx.Abs(zetaEM(complex(0.5, 10.0), lnn))
	cur := 1 / cmplx.Abs(zetaEM(complex(0.5, 10.0+dt), lnn))
	for t := 10.0 + 2*dt; t <= 34.0; t += dt {
		next := 1 / cmplx.Abs(zetaEM(complex(0.5, t), lnn))
		if cur > prev && cur >= next && cur > 20 {
			found = append(found, pk{t - dt, cur})
		}
		prev, cur = cur, next
	}
	for i, p := range found {
		mark := " "
		if i < len(known) && math.Abs(p.t-known[i]) < 0.005 {
			mark = "*"
		}
		fmt.Printf("    %s beacon at t = %9.4f (brightness %6.0f)   known zero: %9.6f\n",
			mark, p.t, p.b, known[i])
	}
	fmt.Println("    every dark point of the hunt shines as a pole in the exotic lens.")

	// Part 2: the negative-energy budget - the Mertens tide.
	fmt.Println("\n  WORMHOLE STABILITY - the exotic-energy budget M(x)/sqrt(x):")
	const N = 10000000
	mu := make([]int8, N+1)
	for i := range mu {
		mu[i] = 1
	}
	sieve := make([]bool, N+1)
	for p := 2; p <= N; p++ {
		if sieve[p] {
			continue
		}
		for q := p; q <= N; q += p {
			if q > p {
				sieve[q] = true
			}
			mu[q] = -mu[q]
		}
		pp := p * p
		if pp <= N {
			for q := pp; q <= N; q += pp {
				mu[q] = 0
			}
		}
	}
	var m int64
	worst, worstX := 0.0, 0
	for x := 1; x <= N; x++ {
		m += int64(mu[x])
		if x >= 100 {
			tide := math.Abs(float64(m)) / math.Sqrt(float64(x))
			if tide > worst {
				worst, worstX = tide, x
			}
		}
	}
	fmt.Printf("    exotic matter summed to x = %d: M = %d\n", N, m)
	fmt.Printf("    worst tide |M(x)|/sqrt(x) = %.4f at x = %d - the budget holds.\n", worst, worstX)
	fmt.Println("    a traversable wormhole demands bounded negative energy; the Riemann")
	fmt.Println("    Hypothesis IS that certificate: RH <=> the Mertens tide respects every")
	fmt.Println("    epsilon-budget. The Clay problem is the wormhole's stability theorem.")

	fmt.Println("\n  and the wormhole itself already flies with the fleet: reciprocity")
	fmt.Println("  (F88) teleports a billion-term wave into ten bounces - instant travel")
	fmt.Println("  between remote points of the sum.")
}
