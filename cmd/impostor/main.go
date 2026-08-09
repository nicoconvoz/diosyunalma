// Command impostor is the second-order decoy the chest demanded.
//
// Findings 57 and 61 excluded pairwise mechanisms for the chest's teeth;
// Finding 63 crowned gap 12's +2.2% as the landscape's one scale-free
// invariant. This decoy settles the level: a synthetic sequence carrying
// ALL the shallow structure of the primes — the exact density drift, the
// full wheel of 2·3·5·7·11·13 (hence every pair and constellation
// correlation those primes generate) — but whose occupancies are otherwise
// INDEPENDENT coin flips, with no deep arithmetic dependence at all.
//
// PRE-REGISTERED: the impostor must show NO teeth — every s₃/s₂² within
// noise of 1, and in particular NOTHING at gap 12 — or the teeth are
// shallow after all. If the impostor is flat while the primes bite, the
// invariant is proven to live above everything an independent-given-wheel
// model can express.
//
// Usage:
//
//	go run ./cmd/impostor [-limit N] [-seed S]
package main

import (
	"flag"
	"fmt"
	"math"
)

const wheel = 30030 // 2*3*5*7*11*13

func main() {
	limit := flag.Int("limit", 100_000_000, "build the impostor up to this value")
	seed := flag.Uint64("seed", 2026, "random seed")
	flag.Parse()

	coprime := make([]bool, wheel)
	phi := 0
	for r := 0; r < wheel; r++ {
		if gcd(r, wheel) == 1 {
			coprime[r] = true
			phi++
		}
	}
	boost := float64(wheel) / float64(phi)

	// splitmix64.
	state := *seed
	next := func() float64 {
		state += 0x9E3779B97F4A7C15
		z := state
		z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
		z = (z ^ (z >> 27)) * 0x94D049BB133111EB
		z ^= z >> 31
		return float64(z>>11) / (1 << 53)
	}

	bits := make([]byte, *limit/8+1)
	members := []int{}
	for m := 100; m < *limit; m++ {
		if !coprime[m%wheel] {
			continue
		}
		if next() < boost/math.Log(float64(m)) {
			bits[m/8] |= 1 << (m % 8)
			members = append(members, m)
		}
	}
	isM := func(x int) bool { return x < *limit && bits[x/8]&(1<<(x%8)) != 0 }
	fmt.Printf("THE IMPOSTOR — %d fake primes on the full wheel of 13 (seed %d)\n", len(members), *seed)
	fmt.Println("\n   d   ratio     deviation    z     real primes at 1e9 said")

	real9 := map[int]string{6: "-0.69%", 12: "+2.20%", 18: "+0.14%", 24: "-0.10%",
		30: "-2.10%", 36: "-12.51%", 42: "-0.03%", 48: "-9.07%"}
	for _, d := range []int{6, 12, 18, 24, 30, 36, 42, 48} {
		var nPair, kPair, nTriple, kTriple float64
		j := 0
		for i, p := range members {
			if p+2*d >= *limit {
				break
			}
			for j < len(members) && members[j] < p+d {
				j++
			}
			if j >= len(members) || members[j] != p+d {
				continue
			}
			nPair++
			firstEmpty := i+1 < len(members) && members[i+1] == p+d
			if firstEmpty {
				kPair++
			}
			if isM(p + 2*d) {
				nTriple++
				if firstEmpty && j+1 < len(members) && members[j+1] == p+2*d {
					kTriple++
				}
			}
		}
		s2 := kPair / nPair
		s3 := kTriple / nTriple
		ratio := s3 / (s2 * s2)
		rel := (1 - s3) / (s3 * nTriple)
		rel += 4 * (1 - s2) / (s2 * nPair)
		sigma := ratio * math.Sqrt(rel)
		fmt.Printf("  %2d  %.4f   %+7.2f%%  %+6.1f     %s\n",
			d, ratio, 100*(ratio-1), (ratio-1)/sigma, real9[d])
	}
	fmt.Println("\nif every impostor row sits at 1 while the primes bite, the teeth -")
	fmt.Println("and the invariant of 12 - live above everything a wheel of coins can say.")
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
