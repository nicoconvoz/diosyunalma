// Command scalpel is the re-sharpened blade Finding 57 demanded.
//
// The naive echo drowned in the universal Hardy–Littlewood pair attraction.
// The fix is a family comparison: for every d in {6, 12, 18, 24, 36, 48, 54}
// the pair singular series is IDENTICAL (only 2 and 3 divide each), so pure
// HL background predicts the SAME cross-corridor attraction for all seven.
// Any spread within the family is the anomaly itself, cut free of the roar.
//
// For each d, the echo strength E(d) pools every same-note joint excess:
//
//	E(d) = [Σ_a nXY(a) − Σ_a exp(a)] / Σ_a exp(a),
//
// the relative excess of "prime at p+a AND at p+d+a" over independence,
// pooled across offsets, with 1/√(Σexp) errors.
//
// PRE-REGISTERED, from the chest's teeth (Finding 54): within the equal-
// background family, E(12) must sit HIGH, E(36) and E(48) LOW, with E(18)
// and E(24) at the family mean — or the cross-corridor mechanism is not
// the chest's key. d = 30, 42, 60 are shown for context only (their 5- and
// 7-boosts change the background).
//
// Usage:
//
//	go run ./cmd/scalpel [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nicoconvoz/numerosprimos/primes"
)

func main() {
	limit := flag.Int("limit", 1_000_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	ps := primes.Sieve(*limit)
	bits := make([]byte, *limit/8+1)
	for _, p := range ps {
		bits[p/8] |= 1 << (p % 8)
	}
	isPrime := func(n int) bool { return bits[n/8]&(1<<(n%8)) != 0 }

	fmt.Printf("THE SCALPEL — echo strength within the equal-background family (x <= %d)\n", *limit)
	fmt.Println("\n   d   E(d) relative excess   sigma     family   tooth (F54)")

	family := map[int]bool{6: true, 12: true, 18: true, 24: true, 36: true, 48: true, 54: true}
	teeth := map[int]string{6: "-0.7%", 12: "+2.2%", 18: "flat", 24: "flat",
		30: "shrinking", 36: "-12.5%", 42: "flat", 48: "-9.1%", 54: "weak -", 60: "weak -"}

	type res struct {
		d        int
		e, sigma float64
	}
	results := []res{}
	for _, d := range []int{6, 12, 18, 24, 30, 36, 42, 48, 54, 60} {
		nX := make([]float64, d)
		nY := make([]float64, d)
		nXY := make([]float64, d)
		var n float64
		for _, p := range ps {
			if p+2*d > *limit {
				break
			}
			if p < 5 || !isPrime(p+d) || !isPrime(p+2*d) {
				continue
			}
			n++
			for a := 2; a < d; a += 2 {
				x := isPrime(p + a)
				y := isPrime(p + d + a)
				if x {
					nX[a]++
				}
				if y {
					nY[a]++
				}
				if x && y {
					nXY[a]++
				}
			}
		}
		var joint, expSum float64
		for a := 2; a < d; a += 2 {
			exp := nX[a] * nY[a] / n
			if exp < 25 {
				continue
			}
			joint += nXY[a]
			expSum += exp
		}
		e := (joint - expSum) / expSum
		sigma := math.Sqrt(joint) / expSum
		tag := "      "
		if family[d] {
			tag = "B=1   "
		}
		fmt.Printf("  %2d      %+.4f            %.4f    %s %s\n", d, e, sigma, tag, teeth[d])
		if family[d] {
			results = append(results, res{d, e, sigma})
		}
	}

	// family mean from the registered neutral members (18, 24, 54).
	var mean, wsum float64
	for _, r := range results {
		if r.d == 18 || r.d == 24 || r.d == 54 {
			w := 1 / (r.sigma * r.sigma)
			mean += r.e * w
			wsum += w
		}
	}
	mean /= wsum
	fmt.Printf("\nfamily baseline (weighted mean of 18, 24, 54): %+.4f\n", mean)
	fmt.Println("\n   d   E(d) - baseline   z vs baseline")
	for _, r := range results {
		z := (r.e - mean) / r.sigma
		fmt.Printf("  %2d      %+.4f          %+6.1f\n", r.d, r.e-mean, z)
	}
	fmt.Println("\nequal singular series, equal background - whatever spreads the family")
	fmt.Println("is the chest's own key, cut free of the Hardy-Littlewood roar.")
}
