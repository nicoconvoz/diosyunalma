// Command tritone is the sentence on √2.
//
// Finding 33 reduced the step-factor mystery to one ratio: E = 2 − c₀/c, so
// E = √2 exactly iff c₀/c → 2−√2 = 0.58579. This run computes where c₀/c is
// actually going, using the standard Hardy–Littlewood gap model (the one
// behind the jumping champions): P(d) ∝ B(d)·e^(−d/λ), with B(d) the
// Hardy–Littlewood weight Π (p−1)/(p−2) over odd primes dividing d, and λ the
// mean gap, which grows without bound.
//
// The model must first EARN its extrapolation: matched to each measured mean
// gap, it has to reproduce the measured c₀/c at 10^6, 10^7 and 10^8. Only
// then is its λ → ∞ limit worth reading.
//
// The limit is also computable by hand. Averaged over d, the p ≥ 5 factors of
// B(d)² are class-independent, so only the factor 2 at the prime 3 survives:
// class 0 carries weight 4 against 1 for each other class, giving
//
//	c₀/c → (1/3·4) / (1/3·4 + 2/3·1) = 2/3      E → 2 − 2/3 = 4/3
//
// Usage:
//
//	go run ./cmd/tritone [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/primes"
)

// hlWeight returns B(d) = Π (p−1)/(p−2) over odd primes p dividing d.
func hlWeight(d int) float64 {
	b := 1.0
	for p := 3; p*p <= d; p += 2 {
		if d%p == 0 {
			b *= float64(p-1) / float64(p-2)
			for d%p == 0 {
				d /= p
			}
		}
	}
	if d > 2 && d%2 == 1 {
		b *= float64(d-1) / float64(d-2)
	}
	return b
}

// modelC0C evaluates the quiet share of collisions in the Hardy–Littlewood
// bag baked at mean-gap scale lambda.
func modelC0C(lambda float64) (c0c, meanGap float64) {
	dMax := int(14 * lambda)
	if dMax < 200 {
		dMax = 200
	}
	var c0, c, wSum, dwSum float64
	for d := 2; d <= dMax; d += 2 {
		w := hlWeight(d) * math.Exp(-float64(d)/lambda)
		wSum += w
		dwSum += float64(d) * w
		p2 := w * w
		c += p2
		if d%3 == 0 {
			c0 += p2
		}
	}
	return c0 / c, dwSum / wSum
}

// lambdaFor finds the lambda whose model mean gap matches the target.
func lambdaFor(target float64) float64 {
	lo, hi := 1.0, 200.0
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		if _, m := modelC0C(mid); m < target {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 1_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e6")
		os.Exit(1)
	}

	walk := primes.From(primes.Sieve(*limit), 5)

	fmt.Println("A) VALIDATION — the model must earn its extrapolation")
	fmt.Printf("%-8s %-11s %-12s %-12s %s\n",
		"N", "mean gap", "c0/c meas.", "c0/c model", "model/meas")

	for _, n := range []int{1_000_000, 10_000_000, 100_000_000, 1_000_000_000} {
		if n > *limit {
			continue
		}
		cut := sort.SearchInts(walk, n+1)
		gaps := primes.Gaps(walk[:cut])

		freq := map[int]int{}
		sum := 0
		for _, g := range gaps {
			freq[g]++
			sum += g
		}
		meanGap := float64(sum) / float64(len(gaps))

		var c0, c float64
		for g, cnt := range freq {
			p := float64(cnt) / float64(len(gaps))
			c += p * p
			if g%3 == 0 {
				c0 += p * p
			}
		}
		measured := c0 / c

		model, _ := modelC0C(lambdaFor(meanGap))
		fmt.Printf("%-8.0e %-11.3f %-12.4f %-12.4f %.4f\n",
			float64(n), meanGap, measured, model, model/measured)
	}

	fmt.Println("\nB) THE BAKE — where the quiet share is heading as the bag widens")
	fmt.Printf("%-10s %-12s %s\n", "lambda", "c0/c", "E = 2 - c0/c")
	for _, l := range []float64{15, 25, 50, 100, 300, 1000, 3000} {
		v, _ := modelC0C(l)
		fmt.Printf("%-10.0f %-12.4f %.4f\n", l, v, 2-v)
	}
	fmt.Printf("%-10s %-12.4f %.4f   <- the hand-computed limit: 2/3 and 4/3\n",
		"infinity", 2.0/3.0, 4.0/3.0)

	fmt.Println("\nC) THE TRITONE CROSSING")
	target := 2 - math.Sqrt2
	lo, hi := 10.0, 3000.0
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		if v, _ := modelC0C(mid); v < target {
			lo = mid
		} else {
			hi = mid
		}
	}
	lambdaStar := (lo + hi) / 2
	fmt.Printf("c0/c = 2-sqrt(2) = %.5f at lambda* ≈ %.1f  ->  N* ~ e^lambda* ≈ 1e%.0f\n",
		target, lambdaStar, lambdaStar/math.Ln10)
	fmt.Println("\nthe quiet share passes THROUGH the tritone and settles at 2/3:")
	fmt.Println("the step factor's destination is 4/3 — the just perfect fourth.")
}
