// Command compose closes the mechanism for R(d) by predicting its
// consecutiveness part instead of merely measuring it.
//
// THE MODEL, Odlyzko–Rubinstein–Wolf in miniature. A pair of primes at gap d
// is consecutive when the interval between them is empty. The interior holds
// exactly d/3 − 1 integers coprime to 6 (for 6 | d), and if each fails to be
// prime independently, the survival must decay geometrically in that count:
//
//	s₂(d) = s₂(6)^(d/3−1)
//
// One calibration at d = 6, zero further freedom. The triple survival should
// then satisfy s₃ ≈ s₂² (two independent intervals), making the
// consecutiveness factor s₃/s₂² ≈ 1 and closing the level:
//
//	R(d) ≈ C·B(d)      with  C = Π (q−3)(q−1)/(q−2)² = 0.81980245
//
// Deviations of s₃/s₂² from 1 measure the correlation between the two
// intervals' emptiness — Hardy–Littlewood correlations the independent model
// cannot carry.
//
// Usage:
//
//	go run ./cmd/compose [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nicoconvoz/numerosprimos/pattern"
	"github.com/nicoconvoz/numerosprimos/primes"
)

// eulerC is Π (q−3)(q−1)/(q−2)² over primes q ≥ 5 (Finding 20).
const eulerC = 0.81980245

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()

	if *limit < 1_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e6")
		os.Exit(1)
	}

	set := primes.NewSet(*limit)
	walk := primes.From(primes.Sieve(*limit), 5)
	gaps := primes.Gaps(walk)
	n := float64(len(gaps) - 1)
	fmt.Printf("primes above 3 up to %d: %d\n", *limit, len(walk))

	freq := map[int]int{}
	for _, g := range gaps {
		freq[g]++
	}
	joint := map[int]int{}
	for i := 1; i < len(gaps); i++ {
		if gaps[i-1] == gaps[i] {
			joint[gaps[i]]++
		}
	}

	ds := []int{6, 12, 18, 24, 30, 36, 42}

	// ---------------- the survival law ----------------
	fmt.Println("\nA) PAIR SURVIVAL — geometric in the coprime candidate count?")
	fmt.Println("   model: s2(d) = s2(6)^(d/3-1), calibrated once at d=6")
	fmt.Printf("\n   %-5s %-12s %-12s %-12s %s\n",
		"d", "candidates", "s2 measured", "s2 predicted", "ratio")

	s2 := map[int]float64{}
	for _, d := range ds {
		s2[d] = float64(freq[d]) / float64(primes.ProgressionPairs(set, d))
	}
	base := s2[6]
	for _, d := range ds {
		cands := d/3 - 1
		pred := math.Pow(base, float64(cands))
		fmt.Printf("   %-5d %-12d %-12.5f %-12.5f %.4f\n",
			d, cands, s2[d], pred, s2[d]/pred)
	}

	// ---------------- interval independence ----------------
	fmt.Println("\nB) TRIPLE SURVIVAL — are the two intervals independent?")
	fmt.Printf("   %-5s %-12s %-12s %-14s %s\n",
		"d", "s3 measured", "s2^2", "s3/s2^2", "reading")
	for _, d := range ds {
		t := primes.ProgressionTriples(set, d)
		if t == 0 || joint[d] == 0 {
			continue
		}
		s3 := float64(joint[d]) / float64(t)
		ratio := s3 / (s2[d] * s2[d])
		reading := "independent"
		if math.Abs(ratio-1) > 0.08 {
			reading = "correlated (HL)"
		}
		fmt.Printf("   %-5d %-12.5f %-12.5f %-14.4f %s\n",
			d, s3, s2[d]*s2[d], ratio, reading)
	}

	// ---------------- the closed level ----------------
	fmt.Println("\nC) THE COMPOSED LEVEL — R(d) against C·B(d), nothing fitted")
	fmt.Printf("   %-5s %-12s %-12s %s\n", "d", "R measured", "C·B(d)", "ratio")
	for _, d := range ds {
		if joint[d] == 0 {
			continue
		}
		g := float64(freq[d])
		r := float64(joint[d]) * n / (g * g)
		pred := eulerC * pattern.SingularBoost(d)
		fmt.Printf("   %-5d %-12.4f %-12.4f %.4f\n", d, r, pred, r/pred)
	}
	fmt.Println("\n   C = 0.81980245 is the Euler product of Finding 20; B(d) the")
	fmt.Println("   singular boost. The consecutiveness factor enters as ~1 via B.")
}
