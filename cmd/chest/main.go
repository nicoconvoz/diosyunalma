// Command chest opens the lock Finding 52 discovered: gap 12 is the only
// gap that likes repeating. This run produces the full signature of the
// anomaly and kills or crowns two candidate keys.
//
// The statistic is Finding 52's s₃/s₂² — the constellation-conditioned
// probability that a gap d repeats — measured for every d ≡ 0 (mod 6) up to
// 120, at BOTH 10^8 and 10^9, with binomial errors.
//
// PRE-REGISTERED:
//
//	H1 (divisor richness): the anomaly tracks τ(d), the divisor count.
//	   Killed if the correlation across d is weak or counterexampled.
//	H2 (pure size): the anomaly is a function of d/ln x only. Killed if
//	   deviations exist that do NOT shrink from 10^8 to 10^9 at fixed d.
//	The residue after both verdicts is the true signature — the teeth of
//	the key, published as data.
//
// Usage:
//
//	go run ./cmd/chest
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

const dMax = 120

type row struct {
	ratio, z float64
	triples  float64
}

func scan(limit int) map[int]row {
	ps := primes.Sieve(limit)
	out := map[int]row{}
	for d := 6; d <= dMax; d += 6 {
		var nPair, kPair, nTriple, kTriple float64
		j := 0
		for i, p := range ps {
			if p+2*d > limit {
				break
			}
			if p < 5 {
				continue
			}
			for j < len(ps) && ps[j] < p+d {
				j++
			}
			if j >= len(ps) || ps[j] != p+d {
				continue
			}
			nPair++
			firstEmpty := ps[i+1] == p+d
			if firstEmpty {
				kPair++
			}
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
		rel := (1 - s3) / (s3 * nTriple)
		rel += 4 * (1 - s2) / (s2 * nPair)
		out[d] = row{ratio, (ratio - 1) / (ratio * math.Sqrt(rel)), nTriple}
	}
	return out
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

func tau(n int) int {
	t := 0
	for k := 1; k <= n; k++ {
		if n%k == 0 {
			t++
		}
	}
	return t
}

func main() {
	fmt.Println("THE CHEST — the full signature of the repetition anomaly")
	small := scan(100_000_000)
	big := scan(1_000_000_000)

	fmt.Println("\n   d  tau   ratio@1e8   ratio@1e9    z@1e9   verdict")
	type point struct{ dev, t float64 }
	pts := []point{}
	stableNonzero := 0
	for d := 6; d <= dMax; d += 6 {
		a, b := small[d], big[d]
		verdict := "independent"
		if math.Abs(b.z) > 5 {
			devA, devB := a.ratio-1, b.ratio-1
			if math.Abs(devB) >= 0.6*math.Abs(devA) {
				verdict = "STABLE anomaly"
				stableNonzero++
			} else {
				verdict = "shrinking (finite-size)"
			}
		}
		fmt.Printf("  %3d  %2d    %.4f      %.4f     %+6.1f   %s\n",
			d, tau(d), a.ratio, b.ratio, b.z, verdict)
		pts = append(pts, point{b.ratio - 1, float64(tau(d))})
	}

	// H1: correlation of the 1e9 deviation with tau(d).
	var sx, sy, sxx, syy, sxy float64
	n := float64(len(pts))
	for _, p := range pts {
		sx += p.t
		sy += p.dev
		sxx += p.t * p.t
		syy += p.dev * p.dev
		sxy += p.t * p.dev
	}
	corr := (n*sxy - sx*sy) / math.Sqrt((n*sxx-sx*sx)*(n*syy-sy*sy))
	fmt.Printf("\nH1 divisor richness: corr(deviation, tau) = %+.3f\n", corr)
	fmt.Printf("H2 pure size: stable anomalies at fixed d = %d (any > 0 kills H2)\n",
		stableNonzero)
}
