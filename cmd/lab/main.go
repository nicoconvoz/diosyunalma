// Command lab runs a detector against shuffled-gap decoys and reports a verdict
// for every window size or lag.
//
// Usage:
//
//	go run ./cmd/lab [-detector palindrome|lag] [-limit N] [-trials N] [-seed N] [-max N]
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/nicoconvoz/diosyunalma/control"
	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

// detector is any measurement that reduces a gap sequence to a single count at
// a given parameter. Adding a new experiment means adding one of these; the
// control harness stays untouched.
type detector struct {
	label string
	param string
	count func(gaps []int, n int) int
}

var detectors = map[string]detector{
	"palindrome": {
		label: "windows of k gaps that read the same in both directions",
		param: "k",
		count: pattern.Palindromes,
	},
	"lag": {
		label: "positions where two gaps separated by d are equal",
		param: "d",
		count: pattern.LagEquality,
	},
}

func main() {
	name := flag.String("detector", "palindrome", "which detector to run: palindrome or lag")
	limit := flag.Int("limit", 1_000_000, "sieve primes up to this value")
	trials := flag.Int("trials", 40, "decoy sequences drawn per parameter value")
	seed := flag.Int64("seed", 2026, "random seed, so runs are reproducible")
	max := flag.Int("max", 9, "largest parameter value to measure")
	flag.Parse()

	det, ok := detectors[*name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown detector %q: use palindrome or lag\n", *name)
		os.Exit(1)
	}
	if *limit < 2 || *trials < 1 || *max < 1 {
		fmt.Fprintln(os.Stderr, "limit must be >= 2, trials >= 1, max >= 1")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))

	seq := primes.Sieve(*limit)
	gaps := primes.Gaps(seq)

	fmt.Printf("detector: %s — %s\n", *name, det.label)
	fmt.Printf("primes <= %d: %d    gaps: %d    decoys each: %d    seed: %d\n\n",
		*limit, len(seq), len(gaps), *trials, *seed)

	fmt.Printf("%-4s %-12s %-14s %-10s %-10s %s\n",
		det.param, "observed", "decoy mean", "ratio", "z", "verdict")

	for n := 1; n <= *max; n++ {
		observed := det.count(gaps, n)
		if observed == 0 && n == 1 && *name == "palindrome" {
			continue // k=1 is not a window
		}

		result := control.Evaluate(observed, *trials, func() int {
			return det.count(control.ShuffleGaps(gaps, rng), n)
		})

		fmt.Printf("%-4d %-12d %-14.1f %-10.3f %-10.1f %s\n",
			n, result.Observed, result.Mean, result.Ratio, result.ZScore, verdict(result))
	}
}

func verdict(r control.Result) string {
	if !r.Significant() {
		return "noise"
	}
	if r.ZScore > 0 {
		return "EXCESS"
	}
	return "DEFICIT"
}
