// Command consecutive isolates the cost of demanding that primes be adjacent.
//
// Open question 1 asks what produces a deficit the Hardy–Littlewood singular
// series does not predict. The standing suspect is the consecutiveness
// requirement: this project measures CONSECUTIVE primes with equal gaps, while
// the k-tuple conjecture counts p, p+d, p+2d all prime regardless of what lies
// between them.
//
// The test needs no new theory. Measure both counts. If the unconstrained one
// tracks the singular series while the consecutive one does not, the suspect is
// convicted. If neither tracks it, the diagnosis was wrong.
//
// Usage:
//
//	go run ./cmd/consecutive [-limit N]
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

func main() {
	limit := flag.Int("limit", 10_000_000, "sieve primes up to this value")
	flag.Parse()

	if *limit < 100 {
		fmt.Fprintln(os.Stderr, "limit must be at least 100")
		os.Exit(1)
	}

	set := primes.NewSet(*limit)
	walk := primes.From(primes.Sieve(*limit), 5)
	gaps := primes.Gaps(walk)

	fmt.Printf("primes <= %d, above 3: %d    gaps: %d\n\n", *limit, len(walk), len(gaps))

	// Consecutive count: two adjacent gaps both equal to d.
	consecutive := map[int]int{}
	for i := 1; i < len(gaps); i++ {
		if gaps[i-1] == gaps[i] {
			consecutive[gaps[i]]++
		}
	}

	fmt.Println("A) UNCONSTRAINED — p, p+d, p+2d all prime (what Hardy-Littlewood predicts)")
	fmt.Printf("   %-6s %-12s %-10s %-12s %s\n", "d", "triples", "boost", "per boost", "vs d=6")

	ds := []int{6, 12, 18, 24, 30, 36, 42, 48, 54, 60}
	triples := map[int]int{}
	var base float64
	for _, d := range ds {
		n := primes.ProgressionTriples(set, d)
		triples[d] = n
		b := pattern.SingularBoost(d)
		perBoost := float64(n) / b
		if d == 6 {
			base = perBoost
		}
		fmt.Printf("   %-6d %-12d %-10.4f %-12.1f %.4f\n", d, n, b, perBoost, perBoost/base)
	}

	fmt.Println("\nB) CONSECUTIVE — the same d twice in a row, nothing in between")
	fmt.Printf("   %-6s %-12s %-14s %-12s %s\n",
		"d", "consecutive", "unconstrained", "survival", "vs d=6")

	var survivalBase float64
	for _, d := range ds {
		if triples[d] == 0 {
			continue
		}
		survival := float64(consecutive[d]) / float64(triples[d])
		if d == 6 {
			survivalBase = survival
		}
		fmt.Printf("   %-6d %-12d %-14d %-12.5f %.4f\n",
			d, consecutive[d], triples[d], survival, survival/survivalBase)
	}

	fmt.Println("\n   survival = fraction of arithmetic triples whose primes happen to be adjacent.")
	fmt.Println("   If this falls as d grows, the emptiness requirement is the cost.")
}
