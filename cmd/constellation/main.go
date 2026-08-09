// Command constellation tests the mechanism proposed for the break in the
// geometric law (Finding 23).
//
// THE THEORY. A decoy window mirrors by luck: each added pair costs one
// independent coincidence, so decoy counts fall geometrically in k. Real
// windows at high k, the theory says, are not coincidences — they are rigid
// admissible prime constellations whose densities follow Hardy–Littlewood and
// decay on a different law. Two different decay laws divided by each other
// produce a ratio that bends, which is what Finding 23 measured.
//
// TWO FALSIFIABLE PREDICTIONS.
//
//  1. CONCENTRATION. If high-k windows are constellations, they cluster on a
//     few tight admissible patterns, repeated. Decoy windows scatter across
//     many patterns, rarely repeating.
//
//  2. SCALING IN N. If two decay laws are being divided, the ratio at fixed
//     high k must move with N. Where coincidences still dominate — low k —
//     the ratio should sit still.
//
// Usage:
//
//	go run ./cmd/constellation [-limit N] [-trials N] [-seed N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/nicoconvoz/numerosprimos/control"
	"github.com/nicoconvoz/numerosprimos/pattern"
	"github.com/nicoconvoz/numerosprimos/primes"
)

func centreFree(w []int) bool { return w[len(w)/2]%3 != 0 }

func key(w []int) string {
	parts := make([]string, len(w))
	for i, v := range w {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

func span(w []int) int {
	s := 0
	for _, g := range w {
		s += g
	}
	return s
}

// admissible reports whether the k+1 positions cut out by the window avoid
// covering every residue class mod q. A pattern that covers all of them can
// contain at most one prime tuple in the whole number line — it is not a
// constellation, it is an accident.
func admissible(w []int, q int) bool {
	seen := map[int]bool{0: true}
	pos := 0
	for _, g := range w {
		pos += g
		seen[pos%q] = true
	}
	return len(seen) < q
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	trials := flag.Int("trials", 30, "decoys per limit in the scaling test")
	seed := flag.Int64("seed", 2026, "random seed")
	flag.Parse()

	if *limit < 1_000_000 || *trials < 5 {
		fmt.Fprintln(os.Stderr, "limit >= 1e6 and trials >= 5")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))
	walk := primes.From(primes.Sieve(*limit), 5)
	gaps := primes.Gaps(walk)
	fmt.Printf("primes above 3 up to %d    gaps: %d\n", *limit, len(gaps))

	// ---------------- PREDICTION 1: CONCENTRATION ----------------
	fmt.Println("\n1) CONCENTRATION — what are the high-k windows made of?")

	for _, k := range []int{9, 11} {
		real := pattern.Windows(gaps, k, centreFree)
		decoy := pattern.Windows(control.ShuffleGaps(gaps, rng), k, centreFree)

		fmt.Printf("\n   k=%d   real: %d windows, %d distinct patterns   decoy: %d windows, %d distinct\n",
			k, len(real), distinct(real), len(decoy), distinct(decoy))

		tally := map[string]int{}
		example := map[string][]int{}
		for _, w := range real {
			tally[key(w)]++
			example[key(w)] = w
		}
		type entry struct {
			pat string
			n   int
		}
		list := []entry{}
		for p, n := range tally {
			list = append(list, entry{p, n})
		}
		sort.Slice(list, func(i, j int) bool { return list[i].n > list[j].n })

		top := 6
		if len(list) < top {
			top = len(list)
		}
		fmt.Printf("   %-28s %-6s %-6s %-10s %s\n", "pattern", "times", "span", "adm mod 5", "adm mod 7")
		for _, e := range list[:top] {
			w := example[e.pat]
			fmt.Printf("   %-28s %-6d %-6d %-10v %v\n",
				e.pat, e.n, span(w), admissible(w, 5), admissible(w, 7))
		}
	}

	// ---------------- PREDICTION 2: SCALING IN N ----------------
	fmt.Println("\n2) SCALING — does the ratio at fixed k move with N?")

	limits := []int{}
	for l := 1_000_000; l <= *limit; l *= 10 {
		limits = append(limits, l)
	}
	ks := []int{5, 7, 9, 11}

	fmt.Printf("\n   %-6s", "k")
	for _, l := range limits {
		fmt.Printf(" %-16s", fmt.Sprintf("N=1e%d", int(math.Round(math.Log10(float64(l))))))
	}
	fmt.Println()

	ratios := map[int][]float64{}
	for _, l := range limits {
		cut := sort.SearchInts(walk, l+1)
		g := primes.Gaps(walk[:cut])

		obs := map[int]int{}
		for _, k := range ks {
			obs[k] = pattern.PalindromesWith(g, k, centreFree)
		}

		sum := map[int]float64{}
		sumSq := map[int]float64{}
		for t := 0; t < *trials; t++ {
			d := control.ShuffleGaps(g, rng)
			for _, k := range ks {
				c := float64(pattern.PalindromesWith(d, k, centreFree))
				sum[k] += c
				sumSq[k] += c * c
			}
		}

		for _, k := range ks {
			mean := sum[k] / float64(*trials)
			if mean == 0 {
				ratios[k] = append(ratios[k], math.NaN())
				continue
			}
			ratios[k] = append(ratios[k], float64(obs[k])/mean)
		}
	}

	for _, k := range ks {
		fmt.Printf("   %-6d", k)
		for _, r := range ratios[k] {
			if math.IsNaN(r) {
				fmt.Printf(" %-16s", "—")
				continue
			}
			fmt.Printf(" %-16.4f", r)
		}
		fmt.Println()
	}

	fmt.Println("\n   flat row  -> coincidences: one decay law, the ratio cancels N")
	fmt.Println("   rising row -> two decay laws divided: the constellation regime")
}

func distinct(windows [][]int) int {
	seen := map[string]bool{}
	for _, w := range windows {
		seen[key(w)] = true
	}
	return len(seen)
}
