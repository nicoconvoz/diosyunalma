// Command decompose asks whether the mechanism closes.
//
// PART A — does the product reproduce the measurement?
//
// The observed rate of two equal consecutive gaps, against a model treating the
// two gaps as independent draws, factors exactly:
//
//	R(d) = C(d)·n / G(d)²  =  [ T(d)·n / P₂(d)² ]  ×  [ s₃(d) / s₂(d)² ]
//	                            arithmetic            consecutiveness
//
// where G and C count consecutive pairs and triples, P₂ and T count the same
// shapes without demanding adjacency, and s₂ = G/P₂, s₃ = C/T are the survival
// rates. Every term is measured, nothing is fitted. If the product tracks R(d),
// the two named mechanisms are the whole story and the closure ratio sits at 1.
//
// PART B — which branch of the law carries the odd excess?
//
// The parity law lets an odd palindrome through by either of two routes: its
// centre gap is not a multiple of 3, or no gap in it is. Splitting the count
// shows whether the excess lives in one branch.
//
// Usage:
//
//	go run ./cmd/decompose [-limit N] [-trials N] [-seed N]
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/control"
	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

func main() {
	limit := flag.Int("limit", 10_000_000, "sieve primes up to this value")
	trials := flag.Int("trials", 40, "decoys per window size in part B")
	seed := flag.Int64("seed", 2026, "random seed")
	flag.Parse()

	if *limit < 1000 || *trials < 1 {
		fmt.Fprintln(os.Stderr, "limit must be >= 1000 and trials >= 1")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))
	set := primes.NewSet(*limit)
	walk := primes.From(primes.Sieve(*limit), 5)
	gaps := primes.Gaps(walk)
	n := float64(len(gaps) - 1)

	fmt.Printf("primes above 3 up to %d: %d    gaps: %d\n", *limit, len(walk), len(gaps))

	// ---------------- PART A ----------------
	fmt.Println("\nA) DOES THE PRODUCT CLOSE?")
	fmt.Println("   R(d) measured against arithmetic x consecutiveness, nothing fitted")
	fmt.Printf("\n   %-5s %-9s %-9s %-11s %-11s %-9s %s\n",
		"d", "R(d)", "arith", "consec", "product", "closure", "verdict")

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
		if c >= 50 {
			ds = append(ds, d)
		}
	}
	sort.Ints(ds)

	for _, d := range ds {
		g := float64(freq[d])  // consecutive pairs with gap d
		c := float64(joint[d]) // consecutive triples
		p2 := float64(primes.ProgressionPairs(set, d))
		tr := float64(primes.ProgressionTriples(set, d))
		if p2 == 0 || tr == 0 || g == 0 {
			continue
		}

		measured := c * n / (g * g)
		arithmetic := tr * n / (p2 * p2)
		s2 := g / p2
		s3 := c / tr
		consecutiveness := s3 / (s2 * s2)
		product := arithmetic * consecutiveness

		closure := product / measured
		verdict := "closes"
		if closure < 0.97 || closure > 1.03 {
			verdict = "OFF"
		}

		fmt.Printf("   %-5d %-9.4f %-9.4f %-11.4f %-11.4f %-9.4f %s\n",
			d, measured, arithmetic, consecutiveness, product, closure, verdict)
	}

	// ---------------- PART B ----------------
	fmt.Println("\nB) WHICH BRANCH OF THE LAW CARRIES THE ODD EXCESS?")
	fmt.Printf("\n   %-4s %-11s %-13s %-13s %-11s %s\n",
		"k", "total", "centre free", "all div by 3", "decoy mean", "ratio")

	centreFree := func(w []int) bool { return w[len(w)/2]%3 != 0 }
	allDivisible := func(w []int) bool {
		for _, g := range w {
			if g%3 != 0 {
				return false
			}
		}
		return true
	}

	for _, k := range []int{3, 5, 7} {
		total := pattern.Palindromes(gaps, k)
		free := pattern.PalindromesWith(gaps, k, centreFree)
		div := pattern.PalindromesWith(gaps, k, allDivisible)

		res := control.Evaluate(total, *trials, func() int {
			return pattern.Palindromes(control.ShuffleGaps(gaps, rng), k)
		})

		fmt.Printf("   %-4d %-11d %-13d %-13d %-11.1f %.4f\n",
			k, total, free, div, res.Mean, res.Ratio)
	}

	fmt.Println("\n   Now the same split, each branch against its own decoy expectation:")
	fmt.Printf("\n   %-4s %-14s %-11s %-11s %-9s %s\n",
		"k", "branch", "observed", "decoy mean", "ratio", "z")

	branches := []struct {
		name string
		keep func([]int) bool
	}{
		{"centre free", centreFree},
		{"all div by 3", allDivisible},
	}

	for _, k := range []int{3, 5, 7} {
		for _, b := range branches {
			observed := pattern.PalindromesWith(gaps, k, b.keep)
			res := control.Evaluate(observed, *trials, func() int {
				return pattern.PalindromesWith(control.ShuffleGaps(gaps, rng), k, b.keep)
			})
			fmt.Printf("   %-4d %-14s %-11d %-11.1f %-9.4f %+.1f\n",
				k, b.name, res.Observed, res.Mean, res.Ratio, res.ZScore)
		}
	}
}
