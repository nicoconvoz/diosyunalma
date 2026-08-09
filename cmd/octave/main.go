// Command octave tests the closed form derived from the Markov walk:
//
//	E = 2 − c₀/c
//
// where c₀ is the collision rate among gaps divisible by 3 and c the total.
// The derivation: a centre-free palindrome's interior swaps the walk state, so
// the two mirror gaps share one available non-zero class; mirroring happens by
// two routes (both-stay or both-jump), and with q₁ = q₂ the sum collapses to
// 2 − c₀/c. The octave (2) is the ceiling; the quiet gaps discount it.
//
// Falsifiable twice: the LEVEL must land near the measured step (~1.40-1.43),
// and the DRIFT across N must track the measured drift.
package main

import (
	"fmt"
	"math/rand"

	"github.com/nicoconvoz/diosyunalma/control"
	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

func centreFree(w []int) bool { return w[len(w)/2]%3 != 0 }

func main() {
	rng := rand.New(rand.NewSource(2026))
	fmt.Printf("%-8s %-10s %-12s %-14s %-14s %s\n",
		"N", "c0/c", "model 2-c0/c", "measured 3->5", "measured 5->7", "model/meas(3->5)")

	for _, limit := range []int{1_000_000, 10_000_000, 100_000_000} {
		gaps := primes.Gaps(primes.From(primes.Sieve(limit), 5))
		n := float64(len(gaps))

		freq := map[int]int{}
		for _, g := range gaps {
			freq[g]++
		}
		c0, c := 0.0, 0.0
		for g, cnt := range freq {
			p := float64(cnt) / n
			c += p * p
			if g%3 == 0 {
				c0 += p * p
			}
		}
		model := 2 - c0/c

		real3 := pattern.PalindromesWith(gaps, 3, centreFree)
		real5 := pattern.PalindromesWith(gaps, 5, centreFree)
		real7 := pattern.PalindromesWith(gaps, 7, centreFree)

		const trials = 60
		var d3, d5, d7 float64
		for t := 0; t < trials; t++ {
			d := control.ShuffleGaps(gaps, rng)
			d3 += float64(pattern.PalindromesWith(d, 3, centreFree))
			d5 += float64(pattern.PalindromesWith(d, 5, centreFree))
			d7 += float64(pattern.PalindromesWith(d, 7, centreFree))
		}
		stepA := (float64(real5) / float64(real3)) / (d5 / d3)
		stepB := (float64(real7) / float64(real5)) / (d7 / d5)

		fmt.Printf("%-8.0e %-10.4f %-12.4f %-14.4f %-14.4f %.4f\n",
			float64(limit), c0/c, model, stepA, stepB, model/stepA)
	}

	fmt.Println("\nsqrt(2) exactly would require c0/c -> 2 - sqrt(2) = 0.5858")
}
