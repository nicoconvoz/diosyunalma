// Command telescope is the longer instrument the variance question demanded.
//
// Resolution in the zero hunt is 4π over the ln-x range: it grows with the
// LOGARITHM of the sieve limit, so each extra decade of primes buys a fixed
// slice of sharpness. Two improvements compound here: the segmented sieve
// reaches 10^10 where the flat sieve runs out of memory, and the sample
// window extends down to x = 10, which is free range no earlier run used.
// Together: resolution 0.78 → 0.61.
//
// PRE-REGISTERED: if the spacing statistics of Finding 28 were GUE plus
// instrument censoring, the sharper instrument must move the variance and the
// small-gap fraction FURTHER toward GUE and find spacings below the old floor
// of 0.398. If they were truly rigid, the numbers stay put.
//
// Usage:
//
//	go run ./cmd/telescope [-limit N] [-max-gamma G]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/primes"
	"github.com/nicoconvoz/diosyunalma/riemann"
	"github.com/nicoconvoz/diosyunalma/spectral"
)

func main() {
	limit := flag.Int("limit", 10_000_000_000, "sieve primes up to this value (segmented)")
	maxGamma := flag.Float64("max-gamma", 250, "hunt zeros up to this height")
	flag.Parse()

	if *limit < 1_000_000 || *maxGamma < 30 {
		fmt.Fprintln(os.Stderr, "limit >= 1e6 and max-gamma >= 30")
		os.Exit(1)
	}

	const du = 0.005
	uMin, uMax := math.Log(10), math.Log(float64(*limit))
	us := []float64{}
	xs := []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	fmt.Printf("samples: %d, x in [10, %d], resolution ≈ %.2f (was 0.78)\n",
		len(us), *limit, 2*2*math.Pi/(uMax-uMin))

	psi := primes.PsiSegmented(*limit, xs)
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
	for f := 5.0; f <= *maxGamma; f += 0.005 {
		freqs = append(freqs, f)
	}
	power := spectral.Periodogram(us, e, freqs)

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

	zeros := []float64{}
	for _, c := range cands {
		if c.q > threshold {
			zeros = append(zeros, c.f)
		}
	}
	sort.Float64s(zeros)
	fmt.Printf("zeros found up to γ=%.0f: %d\n", *maxGamma, len(zeros))

	fmt.Println("\nCENSUS — completeness check against Riemann–von Mangoldt")
	fmt.Printf("%-8s %-10s %-14s %s\n", "T", "found", "N(T)", "difference")
	for _, T := range []float64{50, 100, 150, 200, *maxGamma} {
		if T > *maxGamma {
			continue
		}
		c := 0
		for _, z := range zeros {
			if z <= T {
				c++
			}
		}
		pred := riemann.ZeroCount(T)
		fmt.Printf("%-8.0f %-10d %-14.1f %+.1f\n", T, c, pred, float64(c)-pred)
	}

	fmt.Println("\nSPACINGS — the pre-registered comparison")
	sp := []float64{}
	for i := 1; i < len(zeros); i++ {
		mid := (zeros[i] + zeros[i-1]) / 2
		sp = append(sp, (zeros[i]-zeros[i-1])*math.Log(mid/(2*math.Pi))/(2*math.Pi))
	}
	fmt.Printf("%-26s %-10s %-10s %-10s %s\n", "", "mean s", "var s", "min s", "frac s<0.5")
	fmt.Printf("%-26s %-10.3f %-10.3f %-10.3f %.1f%%\n",
		"TELESCOPE (res 0.61)", mean(sp), variance(sp), minOf(sp), 100*fracBelow(sp, 0.5))
	fmt.Printf("%-26s %-10s %-10s %-10s %s\n",
		"old instrument (0.78)", "0.990", "0.126", "0.398", "5.6%")
	fmt.Printf("%-26s %-10s %-10s %-10s %s\n", "GUE", "1.000", "0.178", "→0", "11.2%")
	fmt.Printf("%-26s %-10s %-10s %-10s %s\n", "Poisson", "1.000", "1.000", "—", "39.3%")
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
