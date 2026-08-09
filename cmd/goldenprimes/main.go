// Command goldenprimes measures the collateral effect: the golden ratio
// enters the primes not as a length but as a question.
//
// The defining equation of φ is x² = x + 1. Inside the arithmetic of a prime
// p, that equation either has a solution or it does not — and by quadratic
// reciprocity on the discriminant 5, it does exactly when p ≡ ±1 (mod 5). So
// every prime is GOLDEN (φ exists in its world) or NON-GOLDEN, and Dirichlet
// makes the tribes equal in the long run.
//
// The measured effect: consecutive primes AVOID repeating their golden
// character. The repulsion itself is the Lemke Oliver–Soundararajan
// phenomenon (2016) at modulus 5, read through the quadratic character; the
// golden framing is what the discriminant makes literal. φ never shows up as
// a number in this laboratory — Finding 36 killed that — but its defining
// property partitions the primes and the ordering of the primes feels the
// partition.
//
// Usage:
//
//	go run ./cmd/goldenprimes [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nicoconvoz/diosyunalma/primes"
)

// golden reports whether x² = x + 1 has a solution modulo p, i.e. whether 5
// is a quadratic residue mod p, i.e. p ≡ ±1 (mod 5).
func golden(p int) bool {
	r := p % 5
	return r == 1 || r == 4
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 1_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e6")
		os.Exit(1)
	}

	// Primes above 5: the wheel's own generator cannot be classified by it.
	walk := primes.From(primes.Sieve(*limit), 7)

	nG := 0
	for _, p := range walk {
		if golden(p) {
			nG++
		}
	}
	total := len(walk)
	fmt.Printf("primes above 5 up to %d: %d\n", *limit, total)
	fmt.Printf("golden (p ≡ ±1 mod 5): %d = %.5f    Dirichlet says 1/2\n",
		nG, float64(nG)/float64(total))

	var gg, gn, ng, nn int
	for i := 0; i+1 < len(walk); i++ {
		a, b := golden(walk[i]), golden(walk[i+1])
		switch {
		case a && b:
			gg++
		case a && !b:
			gn++
		case !a && b:
			ng++
		default:
			nn++
		}
	}
	n := float64(gg + gn + ng + nn)
	pG := float64(nG) / float64(total)
	expect := pG*pG + (1-pG)*(1-pG)
	stay := float64(gg+nn) / n
	z := (stay - expect) / math.Sqrt(expect*(1-expect)/n)

	fmt.Println("\nthe chain of consecutive primes across the golden split:")
	fmt.Printf("  P(golden → golden)         = %.5f\n", float64(gg)/float64(gg+gn))
	fmt.Printf("  P(non-golden → non-golden) = %.5f\n", float64(nn)/float64(nn+ng))
	fmt.Printf("  P(stay) = %.5f    independent = %.5f    z = %+.1f\n", stay, expect, z)

	fmt.Println("\nfor comparison, the mod-3 walk of Finding 13 has P(stay) = 0.43036:")
	fmt.Println("a second wheel, a second repulsion — weaker here, still enormous.")
}
