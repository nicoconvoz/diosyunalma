// Command firstpass reproduces the first-pass probes — Findings 4 through 8,
// 11 and 15 — which until now lived only in throwaway scripts outside the
// repository.
//
// Nothing here is new science. It is the debt paid: every number in the
// findings record regenerates from a clean checkout, the killed hypotheses
// included, because a record that can only re-run its survivors is a sales
// brochure.
//
// Usage:
//
//	go run ./cmd/firstpass [-seed N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"

	"github.com/nicoconvoz/diosyunalma/control"
	"github.com/nicoconvoz/diosyunalma/information"
	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

func main() {
	seed := flag.Int64("seed", 2026, "random seed")
	flag.Parse()
	rng := rand.New(rand.NewSource(*seed))

	million := primes.Sieve(1_000_000)

	// ---------------- Finding 4: the modular wheel ----------------
	fmt.Println("FINDING 4 — residue entropy: which modulus makes primes least random?")
	fmt.Printf("%-6s %-10s %s\n", "m", "primes", "Cramér decoy")
	cramer := control.CramerDecoy(1_000_000, rng)
	for _, m := range []int{2, 3, 4, 6, 10, 12, 30, 47} {
		fmt.Printf("%-6d %-10.4f %.4f\n",
			m, information.ResidueEntropy(million, m), information.ResidueEntropy(cramer, m))
	}

	// ---------------- Finding 5: the emirp kill ----------------
	fmt.Println("\nFINDING 5 — digit-reversal pairs below 100 (the kill, reproduced)")
	hundred := primes.Sieve(100)
	observed := pattern.ReversalPairs(hundred, 100)

	const trials = 2000
	sum, atLeast := 0, 0
	for t := 0; t < trials; t++ {
		c := pattern.ReversalPairs(control.OddDecoy(100, len(hundred), rng), 100)
		sum += c
		if c >= observed {
			atLeast++
		}
	}
	fmt.Printf("observed %d pairs; decoys mean %.2f; p = %.3f  -> %s\n",
		observed, float64(sum)/trials, float64(atLeast)/trials,
		verdictP(float64(atLeast)/trials))

	// ---------------- Finding 6: Gilbreath ----------------
	fmt.Println("\nFINDING 6 — Gilbreath survival vs shuffled gaps (control imperfect)")
	thousand := primes.Sieve(1000)
	ok, _ := GilbreathOnPrimes(thousand)
	survivors := 0
	const gTrials = 500
	gaps := primes.Gaps(thousand)
	for t := 0; t < gTrials; t++ {
		rest := control.ShuffleGaps(gaps[1:], rng)
		seq := []int{2, 3}
		cur := 3
		for _, d := range rest {
			cur += d
			seq = append(seq, cur)
		}
		if _, broke := pattern.GilbreathRows(seq, 60); broke == -1 {
			survivors++
		}
	}
	fmt.Printf("real primes: %d/60 rows; decoys surviving all 60: %d/%d (%.1f%%)\n",
		ok, survivors, gTrials, 100*float64(survivors)/gTrials)

	// ---------------- Finding 7: unfolding ----------------
	fmt.Println("\nFINDING 7 — unfolded spacings vs the exponential law (granularity caveat)")
	s := primes.Unfolded(million)
	mean := 0.0
	for _, v := range s {
		mean += v
	}
	fmt.Printf("mean spacing: %.4f (want ~1)\n", mean/float64(len(s)))
	fmt.Printf("%-14s %-10s %-10s %s\n", "range", "primes", "Poisson", "ratio")
	edges := []float64{0, 0.25, 0.5, 1.0, 2.0, 100}
	for b := 0; b+1 < len(edges); b++ {
		lo, hi := edges[b], edges[b+1]
		c := 0
		for _, v := range s {
			if v >= lo && v < hi {
				c++
			}
		}
		obs := float64(c) / float64(len(s))
		exp := math.Exp(-lo) - math.Exp(-hi)
		fmt.Printf("[%.2f, %.2f)   %-10.4f %-10.4f %.2fx\n", lo, hi, obs, exp, obs/exp)
	}

	// ---------------- Findings 8 and 11: extension ----------------
	fmt.Println("\nFINDINGS 8/11 — extension cost per mirrored pair (geometric, not linear)")
	walk := primes.Gaps(primes.From(million, 5))
	fmt.Printf("%-4s %-12s %-14s %s\n", "k", "real ext", "decoy ext", "ratio")
	const eTrials = 40
	for k := 4; k <= 7; k++ {
		real := ratio(pattern.Palindromes(walk, k), pattern.Palindromes(walk, k-2))
		dsum := 0.0
		for t := 0; t < eTrials; t++ {
			d := control.ShuffleGaps(walk, rng)
			dsum += ratio(pattern.Palindromes(d, k), pattern.Palindromes(d, k-2))
		}
		decoy := dsum / eTrials
		fmt.Printf("%-4d %-12.5f %-14.5f %.3f\n", k, real, decoy, real/decoy)
	}

	// ---------------- Finding 15: form and content ----------------
	fmt.Println("\nFINDING 15 — the shortfall splits into form × content")
	g := primes.Gaps(primes.From(primes.Sieve(10_000_000), 5))

	q := 0
	for _, v := range g {
		if v%3 == 0 {
			q++
		}
	}
	pStay := float64(q) / float64(len(g))

	bothZero, equalBoth := 0, 0
	for i := 1; i < len(g); i++ {
		if g[i-1]%3 == 0 && g[i]%3 == 0 {
			bothZero++
			if g[i-1] == g[i] {
				equalBoth++
			}
		}
	}
	freq := map[int]int{}
	for _, v := range g {
		freq[v]++
	}
	collAll, collZero := 0.0, 0.0
	for v, c := range freq {
		p := float64(c) / float64(len(g))
		collAll += p * p
		if v%3 == 0 {
			collZero += p * p
		}
	}

	form := (float64(bothZero) / float64(len(g)-1)) / (pStay * pStay)
	content := (float64(equalBoth) / float64(bothZero)) / (collZero / (pStay * pStay))
	fmt.Printf("form (transition matrix)  : %.5f\n", form)
	fmt.Printf("content (values)          : %.5f\n", content)
	fmt.Printf("product                   : %.5f   (Finding 2 measured 0.824)\n", form*content)
}

// GilbreathOnPrimes wraps the detector for the real sequence.
func GilbreathOnPrimes(seq []int) (int, int) {
	return pattern.GilbreathRows(seq, 60)
}

func ratio(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func verdictP(p float64) string {
	if p > 0.05 {
		return "not significant; the pattern is dead"
	}
	return "significant"
}
