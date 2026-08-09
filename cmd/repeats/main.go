// Command repeats tests the registered hypothesis of Finding 38: that the 11%
// suppression of double-quietness concentrates in EXACT residue repeats.
//
// The eight residues mod 15 split into four combined classes of two members
// each — (mod-3 lane, golden character) — so a both-stay transition has
// exactly two possible destinations: the SAME residue, or its single PARTNER.
// If Lemke Oliver–Soundararajan avoidance is deepest at exact repetition, the
// diagonal carries the suppression and the partner cells sit near
// independence. The full 8×8 is aggregated into five categories, each scored
// against the independence expectation built from measured marginals.
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

var residues = []int{1, 2, 4, 7, 8, 11, 13, 14}

func golden(r int) bool { m := r % 5; return m == 1 || m == 4 }

func main() {
	walk := primes.From(primes.Sieve(100_000_000), 7)

	idx := map[int]int{}
	for i, r := range residues {
		idx[r] = i
	}

	var counts [8][8]float64
	var marg [8]float64
	for i := 0; i+1 < len(walk); i++ {
		a, b := idx[walk[i]%15], idx[walk[i+1]%15]
		counts[a][b]++
		marg[a]++
	}
	n := float64(len(walk) - 1)
	for i := range marg {
		marg[i] /= n
	}

	category := func(a, b int) string {
		ra, rb := residues[a], residues[b]
		same3 := ra%3 == rb%3
		sameG := golden(ra) == golden(rb)
		switch {
		case a == b:
			return "exact repeat"
		case same3 && sameG:
			return "partner (both stay)"
		case same3:
			return "only mod-3 stays"
		case sameG:
			return "only golden stays"
		default:
			return "both switch"
		}
	}

	obs := map[string]float64{}
	exp := map[string]float64{}
	for a := 0; a < 8; a++ {
		for b := 0; b < 8; b++ {
			cat := category(a, b)
			obs[cat] += counts[a][b]
			exp[cat] += marg[a] * marg[b] * n
		}
	}

	fmt.Printf("primes above 5 up to 1e8: %d transitions\n\n", int(n))
	fmt.Printf("%-22s %-12s %-12s %-9s %s\n", "category", "observed", "independent", "ratio", "z")
	order := []string{"exact repeat", "partner (both stay)", "only mod-3 stays",
		"only golden stays", "both switch"}
	for _, cat := range order {
		z := (obs[cat] - exp[cat]) / math.Sqrt(exp[cat])
		fmt.Printf("%-22s %-12.0f %-12.0f %-9.4f %+.0f\n",
			cat, obs[cat], exp[cat], obs[cat]/exp[cat], z)
	}

	exact := obs["exact repeat"] / exp["exact repeat"]
	partner := obs["partner (both stay)"] / exp["partner (both stay)"]
	fmt.Printf("\nthe split of double-quietness: exact %.4f vs partner %.4f\n", exact, partner)
	switch {
	case partner > 0.98 && exact < 0.85:
		fmt.Println("VERDICT: the suppression lives ENTIRELY in exact repeats - mechanism named.")
	case exact < partner-0.03:
		fmt.Println("VERDICT: exact repeats carry MOST of the suppression; the partner")
		fmt.Println("cells carry a remainder - the mechanism is named but not exclusive.")
	default:
		fmt.Println("VERDICT: suppression is spread evenly - the hypothesis dies.")
	}
}
