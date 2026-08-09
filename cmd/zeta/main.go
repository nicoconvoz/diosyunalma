// Command zeta hunts the first zero of the Riemann zeta function using nothing
// but a sieve and a periodogram.
//
// THE PREY. The explicit formula writes psi(x) = x − Σ_ρ x^ρ/ρ − …, so the
// normalised deviation E(u) = (psi(e^u) − e^u)/e^(u/2) is, up to small smooth
// corrections, a superposition of cosines in u — one per zeta zero, the first
// at γ₁ = 14.134725. If that is true, a periodogram of E must peak there.
//
// PRE-REGISTERED TARGETS, written before the data is looked at:
//
//	γ₁ = 14.134725   γ₂ = 21.022040   γ₃ = 25.010858
//	γ₄ = 30.424876   γ₅ = 32.935062
//
// THE CONTROL. A Cramér decoy — include n with probability 1/ln n, weight ln n
// — has the same density and the same smooth part, but no arithmetic. Its
// spectrum must show no peaks at the zeros. If it did, the peaks would belong
// to the method, not to the primes.
//
// Usage:
//
//	go run ./cmd/zeta [-limit N] [-seed N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/primes"
	"github.com/nicoconvoz/diosyunalma/spectral"
)

var knownZeros = []float64{14.134725, 21.022040, 25.010858, 30.424876, 32.935062, 37.586178, 40.918719, 43.327073, 48.005151, 49.773832}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	seed := flag.Int64("seed", 2026, "random seed for the Cramér control")
	flag.Parse()

	if *limit < 100_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e5")
		os.Exit(1)
	}

	// Sample grid, uniform in u = ln x.
	const du = 0.005
	uMin, uMax := math.Log(100), math.Log(float64(*limit))
	us := []float64{}
	xs := []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	fmt.Printf("samples: %d points, u in [%.2f, %.2f], x in [100, %d]\n",
		len(us), uMin, uMax, *limit)

	// The real series.
	psi := primes.PsiAt(*limit, xs)
	real := normalise(psi, xs, us)

	// The Cramér control: same density, same smooth part, no arithmetic.
	rng := rand.New(rand.NewSource(*seed))
	decoySum := make([]float64, len(xs))
	{
		sum, next := 0.0, 0
		for n := 3; n <= *limit; n++ {
			if rng.Float64() < 1/math.Log(float64(n)) {
				sum += math.Log(float64(n))
			}
			for next < len(xs) && xs[next] <= n {
				decoySum[next] = sum
				next++
			}
		}
		for ; next < len(xs); next++ {
			decoySum[next] = sum
		}
	}
	decoy := normalise(decoySum, xs, us)

	// The spectra.
	freqs := []float64{}
	for f := 5.0; f <= 52.0; f += 0.005 {
		freqs = append(freqs, f)
	}
	realPower := spectral.Periodogram(us, hann(real), freqs)
	decoyPower := spectral.Periodogram(us, hann(decoy), freqs)

	// ---------------- the verdict ----------------
	fmt.Println("\nTOP PEAKS — the primes")
	fmt.Printf("%-4s %-12s %-12s %-14s %s\n", "#", "frequency", "power", "nearest zero", "distance")
	realPeaks := topPeaks(freqs, realPower, 10)
	hits := 0
	for i, p := range realPeaks {
		z, d := nearestZero(p.f)
		mark := ""
		if d < 0.15 {
			mark = "  <- HIT"
			hits++
		}
		fmt.Printf("%-4d %-12.4f %-12.4g γ=%-12.4f %+.4f%s\n", i+1, p.f, p.power, z, p.f-z, mark)
	}

	fmt.Println("\nTOP PEAKS — the Cramér control (must miss the zeros)")
	fmt.Printf("%-4s %-12s %-12s %-14s %s\n", "#", "frequency", "power", "nearest zero", "distance")
	for i, p := range topPeaks(freqs, decoyPower, 3) {
		z, d := nearestZero(p.f)
		mark := ""
		if d < 0.15 {
			mark = "  <- hit (bad)"
		}
		fmt.Printf("%-4d %-12.4f %-12.4g γ=%-12.4f %+.4f%s\n", i+1, p.f, p.power, z, p.f-z, mark)
	}

	// Power at each zero's own frequency, primes against control. The decoy's
	// total power is larger — a density-matched random walk is red noise, its
	// power piled at low frequencies — so the honest comparison is made AT the
	// zeros, not against the decoy's maximum.
	fmt.Println("\nPOWER AT THE ZEROS — primes vs control")
	fmt.Printf("%-12s %-12s %-12s %s\n", "zero", "primes", "control", "ratio")
	for _, z := range knownZeros {
		i := int(math.Round((z - freqs[0]) / (freqs[1] - freqs[0])))
		if i < 0 || i >= len(freqs) {
			continue
		}
		fmt.Printf("γ=%-10.4f %-12.4g %-12.4g %.0fx\n",
			z, realPower[i], decoyPower[i], realPower[i]/decoyPower[i])
	}

	if len(realPeaks) > 0 {
		best := realPeaks[0]
		fmt.Printf("\nMEASURED γ₁ = %.4f    true 14.134725    error %.3f%%\n",
			best.f, 100*math.Abs(best.f-14.134725)/14.134725)
		fmt.Printf("zeros matched among the top peaks: %d of %d\n", hits, len(realPeaks))
	}
}

// normalise turns a cumulative series S(x) into E(u) = (S − x)/√x.
func normalise(s []float64, xs []int, us []float64) []float64 {
	out := make([]float64, len(s))
	for i := range s {
		x := float64(xs[i])
		out[i] = (s[i] - x) / math.Sqrt(x)
	}
	_ = us
	return out
}

// hann applies a Hann window, taming the spectral leakage that the trend at
// the edges of the sampled range would otherwise smear across all frequencies.
func hann(y []float64) []float64 {
	out := make([]float64, len(y))
	n := float64(len(y) - 1)
	for i, v := range y {
		out[i] = v * 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}
	return out
}

type peak struct {
	f     float64
	power float64
}

// topPeaks finds local maxima, refines each by parabolic interpolation on its
// neighbours, and returns the strongest, at least 0.5 apart in frequency.
func topPeaks(freqs, power []float64, n int) []peak {
	cands := []peak{}
	for i := 1; i+1 < len(power); i++ {
		if power[i] > power[i-1] && power[i] > power[i+1] {
			f := freqs[i]
			den := power[i-1] - 2*power[i] + power[i+1]
			if den != 0 {
				f += 0.5 * (power[i-1] - power[i+1]) / den * (freqs[1] - freqs[0])
			}
			cands = append(cands, peak{f, power[i]})
		}
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].power > cands[j].power })

	out := []peak{}
	for _, c := range cands {
		tooClose := false
		for _, kept := range out {
			if math.Abs(kept.f-c.f) < 0.5 {
				tooClose = true
				break
			}
		}
		if !tooClose {
			out = append(out, c)
		}
		if len(out) == n {
			break
		}
	}
	return out
}

func nearestZero(f float64) (float64, float64) {
	best, dist := knownZeros[0], math.Abs(f-knownZeros[0])
	for _, z := range knownZeros[1:] {
		if d := math.Abs(f - z); d < dist {
			best, dist = z, d
		}
	}
	return best, dist
}
