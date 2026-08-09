// Command residue reproduces the mod-3 findings from the repository.
//
// It regenerates Finding 13 (the operator law and the residue transition
// matrix), Finding 16 (the exceptions) and Finding 17 (the singular-series
// boost), all of which previously existed only as throwaway scripts.
//
// Usage:
//
//	go run ./cmd/residue [-limit N]
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/nicoconvoz/numerosprimos/pattern"
	"github.com/nicoconvoz/numerosprimos/primes"
)

// walkFloor is the first prime the mod-3 residue walk is defined for. 2 and 3
// generate the structure and cannot be bound by it.
const walkFloor = 5

func main() {
	limit := flag.Int("limit", 10_000_000, "sieve primes up to this value")
	maxK := flag.Int("max-k", 9, "largest palindrome window to check the law on")
	flag.Parse()

	if *limit < 2 || *maxK < 2 {
		fmt.Fprintln(os.Stderr, "limit must be >= 2 and max-k >= 2")
		os.Exit(1)
	}

	full := primes.Sieve(*limit)
	walk := primes.From(full, walkFloor)
	gaps := primes.Gaps(walk)

	fmt.Printf("primes <= %d: %d    above %d: %d    gaps: %d\n\n",
		*limit, len(full), walkFloor, len(walk), len(gaps))

	// --- Finding 16: the exceptions -----------------------------------
	fmt.Println("EXCEPTIONS — the cases the system cannot hold")
	allGaps := primes.Gaps(full)
	odd := 0
	for i, g := range allGaps {
		if g%2 != 0 {
			odd++
			fmt.Printf("  odd gap: %d -> %d (gap %d)\n", full[i], full[i+1], g)
		}
	}
	broken := 0
	for k := 2; k <= *maxK; k++ {
		broken += pattern.LawViolations(allGaps, k)
	}
	fmt.Printf("  odd gaps in the full sequence      : %d\n", odd)
	fmt.Printf("  law violations including 2 and 3   : %d\n", broken)

	// --- Finding 13: the operator law ---------------------------------
	fmt.Println("\nOPERATOR LAW — palindromic windows above 3")
	fmt.Printf("  %-4s %-14s %s\n", "k", "palindromes", "violations")
	total := 0
	for k := 2; k <= *maxK; k++ {
		v := pattern.LawViolations(gaps, k)
		total += v
		fmt.Printf("  %-4d %-14d %d\n", k, pattern.Palindromes(gaps, k), v)
	}
	fmt.Printf("  total violations above 3           : %d\n", total)

	// --- Finding 13: the transition matrix ----------------------------
	fmt.Println("\nTRANSITION MATRIX — residues mod 3 (0.5 everywhere means no bias)")
	t := primes.TransitionMatrix(walk, 3)
	stay, move := 0, 0
	for from := 1; from <= 2; from++ {
		row := t[from][1] + t[from][2]
		fmt.Printf("  from %d  -> 1: %.5f   -> 2: %.5f   (n=%d)\n",
			from, share(t[from][1], row), share(t[from][2], row), row)
		stay += t[from][from]
		move += row - t[from][from]
	}
	fmt.Printf("  P(stay) = %.5f — consecutive primes %s repeating their class\n",
		share(stay, stay+move),
		pick(share(stay, stay+move) < 0.5, "AVOID", "prefer"))

	// --- Finding 17: the singular series ------------------------------
	fmt.Println("\nSINGULAR SERIES — observed rate of repeated gaps against the boost")
	fmt.Printf("  %-6s %-11s %-12s %-10s %-10s %s\n",
		"d", "observed", "expected", "R(d)", "boost", "R/boost")

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
	ds := []int{}
	for d, c := range joint {
		if c >= 30 {
			ds = append(ds, d)
		}
	}
	sort.Ints(ds)

	pairs := float64(len(gaps) - 1)
	for _, d := range ds {
		p := float64(freq[d]) / float64(len(gaps))
		expected := p * p * pairs
		r := float64(joint[d]) / expected
		b := pattern.SingularBoost(d)
		fmt.Printf("  %-6d %-11d %-12.1f %-10.4f %-10.4f %.4f\n",
			d, joint[d], expected, r, b, r/b)
	}
}

func share(a, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(a) / float64(total)
}

func pick(cond bool, yes, no string) string {
	if cond {
		return yes
	}
	return no
}
