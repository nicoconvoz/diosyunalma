// Command operator hunts the fingerprints of the Hilbert–Pólya operator.
//
// Nobody has constructed the operator in a century of trying. But an operator
// that cannot be seen still leaves measurable fingerprints on its spectrum,
// and this laboratory has the spectrum — measured from a sieve.
//
// FINGERPRINT 1 — DENSITY OF STATES. If the zeros are the energies of the
// Berry–Keating Hamiltonian H = xp, their census must follow the semiclassical
// count N(T) = (T/2π)ln(T/2πe) + 7/8, which is also Riemann–von Mangoldt.
//
// FINGERPRINT 2 — LEVEL REPULSION. Eigenvalues of a Hermitian operator repel:
// unfolded spacings follow GUE (Wigner surmise ~ s²e^(−4s²/π)), with few small
// gaps and variance ≈ 0.18. Independent random energies follow Poisson: many
// small gaps, variance 1. Montgomery–Odlyzko saw GUE in the zeros; this run
// asks whether OUR measured zeros carry the same fingerprint.
//
// HONEST LIMIT, stated up front: the spectral resolution here is ~0.8 in γ,
// so pairs of zeros closer than that would merge and could fake repulsion.
// The density test guards against it — heavy merging would push the census
// visibly below N(T).
//
// Usage:
//
//	go run ./cmd/operator [-limit N] [-max-gamma G]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/primes"
	"github.com/nicoconvoz/diosyunalma/riemann"
	"github.com/nicoconvoz/diosyunalma/spectral"
)

func main() {
	limit := flag.Int("limit", 1_000_000_000, "sieve primes up to this value")
	maxGamma := flag.Float64("max-gamma", 130, "hunt zeros up to this height")
	seed := flag.Int64("seed", 2026, "seed for the Poisson control")
	flag.Parse()

	if *limit < 1_000_000 || *maxGamma < 30 {
		fmt.Fprintln(os.Stderr, "limit >= 1e6 and max-gamma >= 30")
		os.Exit(1)
	}

	// The series, as in the zeta hunt.
	const du = 0.005
	uMin, uMax := math.Log(100), math.Log(float64(*limit))
	us := []float64{}
	xs := []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	fmt.Printf("samples: %d, resolution in γ ≈ %.2f (Hann-widened)\n",
		len(us), 2*2*math.Pi/(uMax-uMin))

	psi := primes.PsiAt(*limit, xs)
	e := make([]float64, len(psi))
	for i := range psi {
		x := float64(xs[i])
		e[i] = (psi[i] - x) / math.Sqrt(x)
	}
	y := hann(e)

	freqs := []float64{}
	for f := 5.0; f <= *maxGamma; f += 0.005 {
		freqs = append(freqs, f)
	}
	power := spectral.Periodogram(us, y, freqs)

	// Candidate zeros: local maxima, scored by power·γ² (the explicit formula
	// weights each zero by ~2/|ρ|, so power falls as γ⁻² and the compensated
	// score is flat for true zeros).
	type cand struct{ f, q float64 }
	cands := []cand{}
	for i := 1; i+1 < len(power); i++ {
		if power[i] > power[i-1] && power[i] > power[i+1] {
			f := freqs[i]
			den := power[i-1] - 2*power[i] + power[i+1]
			if den != 0 {
				f += 0.5 * (power[i-1] - power[i+1]) / den * 0.005
			}
			cands = append(cands, cand{f, power[i] * f * f})
		}
	}

	// Calibrate the acceptance threshold on the ten zeros already confirmed.
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
	fmt.Printf("threshold: median confirmed score %.1f, accepting q > %.1f\n",
		scores[len(scores)/2], threshold)

	zeros := []float64{}
	for _, c := range cands {
		if c.q > threshold {
			zeros = append(zeros, c.f)
		}
	}
	sort.Float64s(zeros)
	fmt.Printf("zeros found up to γ=%.0f: %d\n", *maxGamma, len(zeros))

	// ---------------- FINGERPRINT 1: density of states ----------------
	fmt.Println("\nFINGERPRINT 1 — density of states (Berry–Keating semiclassical count)")
	fmt.Printf("%-8s %-10s %-14s %s\n", "T", "found", "N(T) predicted", "difference")
	for _, T := range []float64{50, 80, 100, 120, *maxGamma} {
		if T > *maxGamma {
			continue
		}
		n := 0
		for _, z := range zeros {
			if z <= T {
				n++
			}
		}
		pred := riemann.ZeroCount(T)
		fmt.Printf("%-8.0f %-10d %-14.1f %+.1f\n", T, n, pred, float64(n)-pred)
	}

	// ---------------- FINGERPRINT 2: level repulsion ----------------
	fmt.Println("\nFINGERPRINT 2 — unfolded spacings against GUE and Poisson")

	spacings := unfold(zeros)
	mGUE, vGUE, fGUE := wignerGUE()

	fmt.Printf("%-22s %-10s %-10s %-10s %s\n", "", "mean s", "var s", "min s", "frac s<0.5")
	fmt.Printf("%-22s %-10.3f %-10.3f %-10.3f %.1f%%\n",
		"MEASURED ZEROS", mean(spacings), variance(spacings), minOf(spacings),
		100*fracBelow(spacings, 0.5))
	fmt.Printf("%-22s %-10.3f %-10.3f %-10s %.1f%%\n",
		"GUE (operator)", mGUE, vGUE, "→0 repel", 100*fGUE)
	fmt.Printf("%-22s %-10.3f %-10.3f %-10s %.1f%%\n",
		"Poisson (random)", 1.0, 1.0, "no repel", 100*(1-math.Exp(-0.5)))

	// A Poisson control with the same density, unfolded identically.
	rng := rand.New(rand.NewSource(*seed))
	fake := []float64{14.13}
	for fake[len(fake)-1] < *maxGamma {
		g := fake[len(fake)-1]
		meanGap := 2 * math.Pi / math.Log(g/(2*math.Pi))
		fake = append(fake, g+rng.ExpFloat64()*meanGap)
	}
	fs := unfold(fake[:len(fake)-1])
	fmt.Printf("%-22s %-10.3f %-10.3f %-10.3f %.1f%%\n",
		"Poisson control", mean(fs), variance(fs), minOf(fs), 100*fracBelow(fs, 0.5))

	fmt.Println("\ncaveat: resolution ~0.8 merges closer pairs; the census above bounds the loss.")
}

// unfold converts raw gaps to unit-mean spacings using the local zero density
// ln(γ/2π)/2π from the counting function.
func unfold(zeros []float64) []float64 {
	out := []float64{}
	for i := 1; i < len(zeros); i++ {
		mid := (zeros[i] + zeros[i-1]) / 2
		out = append(out, (zeros[i]-zeros[i-1])*math.Log(mid/(2*math.Pi))/(2*math.Pi))
	}
	return out
}

// wignerGUE integrates the GUE surmise P(s) = (32/π²)s²e^(−4s²/π) numerically
// for its mean, variance and small-gap fraction.
func wignerGUE() (mean, variance, fracHalf float64) {
	const h = 1e-4
	var m0, m1, m2, below float64
	for s := 0.0; s < 8; s += h {
		p := 32 / (math.Pi * math.Pi) * s * s * math.Exp(-4*s*s/math.Pi)
		m0 += p * h
		m1 += s * p * h
		m2 += s * s * p * h
		if s < 0.5 {
			below += p * h
		}
	}
	return m1 / m0, m2/m0 - (m1/m0)*(m1/m0), below / m0
}

func hann(y []float64) []float64 {
	out := make([]float64, len(y))
	n := float64(len(y) - 1)
	for i, v := range y {
		out[i] = v * 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}
	return out
}

func mean(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func variance(v []float64) float64 {
	m := mean(v)
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return s / float64(len(v))
}

func minOf(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x < m {
			m = x
		}
	}
	return m
}

func fracBelow(v []float64, cut float64) float64 {
	n := 0
	for _, x := range v {
		if x < cut {
			n++
		}
	}
	return float64(n) / float64(len(v))
}
