// Command flanks resolves the mystery Finding 32 left open.
//
// For a prime triple p, p+d, p+2d to be consecutive, both flanking intervals
// must be empty of primes. Finding 32 measured the independence ratio
// s₃/s₂² ≈ 1 to about 2% at 10^8 — except at d = 30 (0.939) and d = 36
// (0.879), flagged as candidate interval correlations at roughly 2 sigma.
//
// This run repeats the measurement at 10^9 with ten times the triples and
// proper error bars.
//
// PRE-REGISTERED: if the d = 30 and d = 36 departures are genuine interval
// correlations they must persist with |z| > 5 at 10^9; if they were
// finite-sample noise they must regress toward 1 within 2 sigma.
//
// Usage:
//
//	go run ./cmd/flanks [-limit N]
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
	fmt.Printf("THE FLANKS — interval independence at %d (%d primes)\n\n", *limit, len(ps))
	fmt.Println("   d   pairs      s2       triples   s3       s3/s2^2   z vs 1")

	for d := 6; d <= 60; d += 6 {
		var nPair, kPair, nTriple, kTriple float64
		j := 0
		for i, p := range ps {
			if p < 5 || p+2*d > *limit {
				if p+2*d > *limit {
					break
				}
				continue
			}
			// advance j to the position of p+d, if prime.
			for j < len(ps) && ps[j] < p+d {
				j++
			}
			if j >= len(ps) || ps[j] != p+d {
				continue
			}
			// pair (p, p+d): both prime. empty interval iff next prime is p+d.
			nPair++
			firstEmpty := ps[i+1] == p+d
			if firstEmpty {
				kPair++
			}
			// triple: p+2d prime too. both intervals empty?
			if bsearch(ps, p+2*d) {
				nTriple++
				if firstEmpty && j+1 < len(ps) && ps[j+1] == p+2*d {
					kTriple++
				}
			}
		}
		s2 := kPair / nPair
		s3 := kTriple / nTriple
		ratio := s3 / (s2 * s2)
		// var(r)/r^2 = var(s3)/s3^2 + 4 var(s2)/s2^2 (binomial errors).
		rel := s3 * (1 - s3) / (s3 * s3 * nTriple)
		rel += 4 * s2 * (1 - s2) / (s2 * s2 * nPair)
		sigma := ratio * math.Sqrt(rel)
		fmt.Printf("  %2d  %8.0f  %.4f  %8.0f  %.4f   %.4f   %+.1f\n",
			d, nPair, s2, nTriple, s3, ratio, (ratio-1)/sigma)
	}
	fmt.Println("\nFinding 32 at 10^8 read: d=30 -> 0.9392, d=36 -> 0.8787 (~2 sigma).")
	fmt.Println("Real correlations persist with |z| > 5; noise regresses toward 1.")
}

func bsearch(ps []int, x int) bool {
	lo, hi := 0, len(ps)
	for lo < hi {
		mid := (lo + hi) / 2
		if ps[mid] < x {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo < len(ps) && ps[lo] == x
}
