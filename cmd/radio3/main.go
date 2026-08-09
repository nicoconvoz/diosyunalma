// Command radio3 tunes the tribes' own stations.
//
// Radio 1 heard the stations of ALL the primes: the zeta zeros. But each
// residue tribe has its own mother function — a Dirichlet L-function — and
// the IMBALANCE between the golden and non-golden tribes oscillates at the
// zeros of L(s, χ₅), where χ₅ is the quadratic character mod 5.
//
// The instrument is the character-weighted Chebyshev sum
//
//	ψ(x, χ) = Σ χ(p^k)·ln p     over prime powers p^k ≤ x
//
// which for a non-principal character has NO smooth main term: the whole
// signal is the wave. E(u) = ψ(e^u, χ)/e^(u/2), Hann window, periodogram.
//
// PRE-REGISTERED, before looking:
//  1. the zeta stations (14.1347, 21.0220, …) must be ABSENT from this dial;
//  2. the peaks must be stable between 10^7 and 10^8 (same physics, longer
//     antenna);
//  3. the peaks are predicted to be the zeros of L(s, χ₅), checkable against
//     published tables after measurement.
//
// Usage:
//
//	go run ./cmd/radio3 [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/numerosprimos/primes"
	"github.com/nicoconvoz/numerosprimos/spectral"
)

// chi is the quadratic character mod 5: +1 on the golden residues, −1 on the
// others, 0 on multiples of 5.
func chi(m int) float64 {
	switch m % 5 {
	case 1, 4:
		return 1
	case 2, 3:
		return -1
	}
	return 0
}

// psiChi evaluates the character-weighted Chebyshev sum at ascending xs.
func psiChi(limit int, xs []int) []float64 {
	ps := primes.Sieve(limit)

	type event struct {
		at int
		w  float64
	}
	powers := []event{}
	for _, p := range ps {
		if p*p > limit {
			break
		}
		lg := math.Log(float64(p))
		for pk := p * p; pk <= limit; pk *= p {
			powers = append(powers, event{pk, chi(pk) * lg})
		}
	}
	sort.Slice(powers, func(i, j int) bool { return powers[i].at < powers[j].at })

	out := make([]float64, len(xs))
	sum := 0.0
	pi, wi := 0, 0
	for i, x := range xs {
		for pi < len(ps) && ps[pi] <= x {
			sum += chi(ps[pi]) * math.Log(float64(ps[pi]))
			pi++
		}
		for wi < len(powers) && powers[wi].at <= x {
			sum += powers[wi].w
			wi++
		}
		out[i] = sum
	}
	return out
}

func spectrum(limit int) ([]float64, []float64) {
	const du = 0.005
	uMin, uMax := math.Log(100), math.Log(float64(limit))
	us, xs := []float64{}, []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	psi := psiChi(limit, xs)
	e := make([]float64, len(psi))
	for i := range psi {
		e[i] = psi[i] / math.Sqrt(float64(xs[i]))
	}
	n := float64(len(e) - 1)
	for i := range e {
		e[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}
	freqs := []float64{}
	for f := 2.0; f <= 30.0; f += 0.005 {
		freqs = append(freqs, f)
	}
	return freqs, spectral.Periodogram(us, e, freqs)
}

type peak struct{ f, p float64 }

func topPeaks(freqs, power []float64, n int) []peak {
	cands := []peak{}
	for i := 1; i+1 < len(power); i++ {
		if power[i] > power[i-1] && power[i] > power[i+1] {
			f := freqs[i]
			den := power[i-1] - 2*power[i] + power[i+1]
			if den != 0 {
				f += 0.5 * (power[i-1] - power[i+1]) / den * 0.005
			}
			cands = append(cands, peak{f, power[i]})
		}
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].p > cands[j].p })
	out := []peak{}
	for _, c := range cands {
		ok := true
		for _, k := range out {
			if math.Abs(k.f-c.f) < 0.5 {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, c)
		}
		if len(out) == n {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].f < out[j].f })
	return out
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	fmt.Println("RADIO 3 — the golden tribe's own stations (L(s, χ₅))")

	fA, pA := spectrum(*limit)
	fB, pB := spectrum(*limit / 10)
	peaksA := topPeaks(fA, pA, 8)
	peaksB := topPeaks(fB, pB, 8)

	fmt.Printf("\n%-14s %-12s %-16s %s\n",
		"peak (long)", "power", "nearest (short)", "stable?")
	zeta := []float64{14.1347, 21.0220, 25.0109}
	zetaHit := false
	for _, a := range peaksA {
		best, dist := 0.0, math.Inf(1)
		for _, b := range peaksB {
			if d := math.Abs(a.f - b.f); d < dist {
				dist, best = d, b.f
			}
		}
		stable := "yes"
		if dist > 0.15 {
			stable = "NO"
		}
		for _, z := range zeta {
			if math.Abs(a.f-z) < 0.3 {
				zetaHit = true
			}
		}
		fmt.Printf("%-14.4f %-12.4g %-16.4f %s\n", a.f, a.p, best, stable)
	}

	fmt.Printf("\nzeta stations on this dial: %v   (pre-registered: absent)\n", zetaHit)
	fmt.Println("predicted identity of the stable peaks: the zeros of L(s, χ₅) —")
	fmt.Println("check against published tables (LMFDB, modulus 5, quadratic character).")
}
