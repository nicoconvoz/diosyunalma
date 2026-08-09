// Command sundial closes the circle: it reconstructs the primes from this
// laboratory's own measured zeros.
//
// Every experiment so far ran one direction — primes in, zeros out. The
// explicit formula asserts a two-way duality, and only the return journey
// proves it. Each zero is a clock hand rotating at speed γ as u = ln x
// advances; the truncated density 1 − (2/√x)·Σ cos(γ·ln x) must spike where
// the hands align, and the explicit formula says they align exactly at the
// prime powers.
//
// PRE-REGISTERED: the top peaks of the clock, built from measured zeros and
// nothing else, land on prime powers. The control — fake hands at Poisson
// heights with the correct density — points at nothing in particular.
//
// Usage:
//
//	go run ./cmd/sundial [-limit N] [-max-gamma G] [-max-x X]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/nicoconvoz/numerosprimos/primes"
	"github.com/nicoconvoz/numerosprimos/riemann"
	"github.com/nicoconvoz/numerosprimos/spectral"
)

func main() {
	limit := flag.Int("limit", 1_000_000_000, "sieve primes up to this value for the zero hunt")
	maxGamma := flag.Float64("max-gamma", 130, "use zeros up to this height as hands")
	maxX := flag.Float64("max-x", 62, "reconstruct the number line up to here")
	seed := flag.Int64("seed", 2026, "seed for the fake-hands control")
	flag.Parse()

	if *limit < 1_000_000 || *maxGamma < 30 || *maxX < 20 {
		fmt.Fprintln(os.Stderr, "limit >= 1e6, max-gamma >= 30, max-x >= 20")
		os.Exit(1)
	}

	// ---- Step 1: measure the hands, exactly as cmd/operator does ----
	zeros := measureZeros(*limit, *maxGamma)
	fmt.Printf("hands: %d zeros measured up to γ=%.0f\n", len(zeros), *maxGamma)

	// ---- Step 2: read the clock ----
	targets := primePowersUpTo(int(*maxX))
	found := clockPeaks(zeros, *maxX, len(targets))

	fmt.Printf("\nTHE CLOCK READ — top %d alignments vs the %d prime powers ≤ %.0f\n",
		len(found), len(targets), *maxX)
	fmt.Printf("%-4s %-10s %-14s %s\n", "#", "peak at", "nearest target", "verdict")

	hits := 0
	for i, p := range found {
		nearest, dist := nearestOf(p, targets)
		mark := "miss"
		if dist <= 0.35 {
			mark = "HIT"
			hits++
		}
		fmt.Printf("%-4d %-10.2f %-14.0f %s\n", i+1, p, nearest, mark)
	}
	fmt.Printf("\nmeasured zeros : %d of %d alignments on prime powers\n", hits, len(found))

	// ---- Step 3: the control ----
	rng := rand.New(rand.NewSource(*seed))
	fake := []float64{14.13}
	for fake[len(fake)-1] < *maxGamma {
		g := fake[len(fake)-1]
		fake = append(fake, g+rng.ExpFloat64()*2*math.Pi/math.Log(g/(2*math.Pi)))
	}
	fake = fake[:len(fake)-1]

	fakeHits := 0
	for _, p := range clockPeaks(fake, *maxX, len(targets)) {
		if _, dist := nearestOf(p, targets); dist <= 0.35 {
			fakeHits++
		}
	}
	fmt.Printf("fake hands     : %d of %d  (%d Poisson heights, same density)\n",
		fakeHits, len(targets), len(fake))
}

// measureZeros repeats the calibrated peak hunt of cmd/operator.
func measureZeros(limit int, maxGamma float64) []float64 {
	const du = 0.005
	uMin, uMax := math.Log(100), math.Log(float64(limit))
	us := []float64{}
	xs := []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	psi := primes.PsiAt(limit, xs)
	e := make([]float64, len(psi))
	for i := range psi {
		x := float64(xs[i])
		e[i] = (psi[i] - x) / math.Sqrt(x)
	}
	n := float64(len(e) - 1)
	for i := range e {
		e[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}

	freqs := []float64{}
	for f := 5.0; f <= maxGamma; f += 0.005 {
		freqs = append(freqs, f)
	}
	power := spectral.Periodogram(us, e, freqs)

	type cand struct{ f, q float64 }
	cands := []cand{}
	for i := 1; i+1 < len(power); i++ {
		if power[i] > power[i-1] && power[i] > power[i+1] {
			cands = append(cands, cand{freqs[i], power[i] * freqs[i] * freqs[i]})
		}
	}
	confirmed := []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
		37.5872, 40.9264, 43.3211, 48.0105, 49.7752}
	scores := []float64{}
	for _, z := range confirmed {
		best := 0.0
		for _, c := range cands {
			if math.Abs(c.f-z) < 0.2 && c.q > best {
				best = c.q
			}
		}
		if best > 0 {
			scores = append(scores, best)
		}
	}
	sort.Float64s(scores)
	threshold := scores[len(scores)/2] / 4

	out := []float64{}
	for _, c := range cands {
		if c.q > threshold {
			out = append(out, c.f)
		}
	}
	sort.Float64s(out)
	return out
}

// clockPeaks returns the positions of the n strongest local maxima of the
// prime clock over (1, maxX], at least 0.8 apart.
func clockPeaks(gammas []float64, maxX float64, n int) []float64 {
	const dx = 0.01
	type pk struct{ x, v float64 }
	cands := []pk{}
	prev2, prev1 := 0.0, 0.0
	for x := 1.5; x <= maxX; x += dx {
		v := riemann.PrimeClock(gammas, x)
		if prev1 > prev2 && prev1 > v && x-dx > 1.5 {
			cands = append(cands, pk{x - dx, prev1})
		}
		prev2, prev1 = prev1, v
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].v > cands[j].v })

	out := []float64{}
	for _, c := range cands {
		ok := true
		for _, kept := range out {
			if math.Abs(kept-c.x) < 0.8 {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, c.x)
		}
		if len(out) == n {
			break
		}
	}
	sort.Float64s(out)
	return out
}

func primePowersUpTo(limit int) []float64 {
	out := []float64{}
	for _, p := range primes.Sieve(limit) {
		for pk := p; pk <= limit; pk *= p {
			out = append(out, float64(pk))
		}
	}
	sort.Float64s(out)
	return out
}

func nearestOf(x float64, targets []float64) (float64, float64) {
	best, dist := targets[0], math.Abs(x-targets[0])
	for _, t := range targets[1:] {
		if d := math.Abs(x - t); d < dist {
			best, dist = t, d
		}
	}
	return best, dist
}
