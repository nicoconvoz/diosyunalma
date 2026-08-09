// Command budget states every mechanism in one currency: bits per prime.
//
// A finding that says "primes avoid repeating gaps" and a finding that says
// "the residue chain remembers three steps" cannot be compared as they stand.
// Priced in bits they can, and a mechanism that recovers none can be dropped.
//
// Two chains are measured. The gap sequence carries the whole description cost,
// so its entropy is literally bits per prime. The flip sequence keeps only
// whether each gap is divisible by 3 — the form — and its binary alphabet lets
// the memory be measured far deeper before the estimator runs out of data.
//
// Every number is paired with a shuffled control. Conditional entropy estimated
// from a finite sample always drifts downward as the history grows, so the
// control is not decoration: it is the zero line. Real memory is only what
// exceeds it.
//
// Usage:
//
//	go run ./cmd/budget [-limit N] [-seed N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/nicoconvoz/numerosprimos/information"
	"github.com/nicoconvoz/numerosprimos/primes"
)

func main() {
	limit := flag.Int("limit", 10_000_000, "sieve primes up to this value")
	seed := flag.Int64("seed", 2026, "random seed for the shuffled control")
	flag.Parse()

	if *limit < 1000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1000")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))
	walk := primes.From(primes.Sieve(*limit), 5)
	gaps := primes.Gaps(walk)
	flips := primes.Flips(gaps, 3)

	decoyGaps := shuffle(gaps, rng)
	decoyFlips := primes.Flips(decoyGaps, 3)

	fmt.Printf("primes above 3 up to %d: %d    gaps: %d    seed: %d\n",
		*limit, len(walk), len(gaps), *seed)

	// --- the starting line ---------------------------------------------
	noStructure := logBinomial(*limit, len(walk)) / float64(len(walk))
	naive := math.Log2(float64(*limit))

	fmt.Println("\nWHERE THE BUDGET STARTS")
	fmt.Printf("  listing each prime outright        : %7.4f bits/prime\n", naive)
	fmt.Printf("  naming k positions out of N        : %7.4f bits/prime   log2(C(N,k))/k\n", noStructure)

	// --- the gap chain: the real description cost ----------------------
	fmt.Println("\nTHE GAP CHAIN — this is the description cost itself")
	fmt.Printf("  %-26s %-12s %-12s %s\n", "model", "bits/prime", "control", "recovered")

	gapBase := information.ConditionalEntropy(gaps, 0)
	fmt.Printf("  %-26s %-12.4f %-12s %.4f\n",
		"gaps, no memory", gapBase, "—", noStructure-gapBase)

	for order := 1; order <= 2; order++ {
		real := information.ConditionalEntropy(gaps, order)
		ctrl := information.ConditionalEntropy(decoyGaps, order)
		bias := gapBase - ctrl
		honest := gapBase - real - bias

		fmt.Printf("  %-26s %-12.4f %-12.4f %.4f\n",
			fmt.Sprintf("+ memory of %d gap(s)", order), real, ctrl, honest)
	}

	// --- the flip chain: the form, measured deep -----------------------
	fmt.Println("\nTHE FLIP CHAIN — only whether each gap is divisible by 3")
	fmt.Printf("  %-14s %-12s %-12s %-12s %s\n",
		"history", "H(bits)", "control", "raw drop", "real drop")

	flipBase := information.ConditionalEntropy(flips, 0)
	fmt.Printf("  %-14s %-12.6f %-12s %-12s %s\n", "none", flipBase, "—", "—", "—")

	for order := 1; order <= 8; order++ {
		real := information.ConditionalEntropy(flips, order)
		ctrl := information.ConditionalEntropy(decoyFlips, order)

		rawDrop := flipBase - real
		bias := flipBase - ctrl
		fmt.Printf("  %-14d %-12.6f %-12.6f %-12.6f %.6f\n",
			order, real, ctrl, rawDrop, rawDrop-bias)
	}

	fmt.Println("\n  real drop = mutual information the past truly carries, with the")
	fmt.Println("  estimator's own bias measured on the control and subtracted.")
}

func shuffle(v []int, rng *rand.Rand) []int {
	out := append([]int(nil), v...)
	rng.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

// logBinomial returns log2(C(n, k)) — the cost of naming k positions among n
// with no structure available to exploit.
func logBinomial(n, k int) float64 {
	lg := func(x float64) float64 { v, _ := math.Lgamma(x); return v }
	return (lg(float64(n)+1) - lg(float64(k)+1) - lg(float64(n-k)+1)) / math.Ln2
}
