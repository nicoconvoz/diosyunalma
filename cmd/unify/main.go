// Command unify tests whether the two branches of the parity law are one number.
//
// Finding 21 split odd palindromes into the two routes the law allows: a centre
// gap not divisible by 3, or no gap divisible by anything but 3. Both ratios
// move monotonically with k and they move opposite ways:
//
//	centre free    1.451 -> 2.052 -> 2.921
//	all div by 3   0.822 -> 0.580 -> 0.391
//
// Successive factors are near 1.414 and 0.707. If those are sqrt(2) and its
// reciprocal, one branch gains exactly what the other loses, and the product of
// the two ratios is a constant that does not depend on k.
//
// A constant product would mean the law does not create structure — it moves a
// fixed amount of it from one branch to the other, and every ratio in Finding
// 21 collapses to one number plus a scaling rule.
//
// Usage:
//
//	go run ./cmd/unify [-limit N] [-trials N] [-max-k N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/nicoconvoz/diosyunalma/control"
	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

// row holds one window size and how each branch scored against its decoys.
type row struct {
	k         int
	free, div float64
	product   float64
}

func centreFree(w []int) bool { return w[len(w)/2]%3 != 0 }

func allDivisible(w []int) bool {
	for _, g := range w {
		if g%3 != 0 {
			return false
		}
	}
	return true
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	trials := flag.Int("trials", 20, "decoys per branch")
	maxK := flag.Int("max-k", 11, "largest odd window to measure")
	seed := flag.Int64("seed", 2026, "random seed")
	flag.Parse()

	if *limit < 1000 || *trials < 2 || *maxK < 3 {
		fmt.Fprintln(os.Stderr, "limit >= 1000, trials >= 2, max-k >= 3")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))
	gaps := primes.Gaps(primes.From(primes.Sieve(*limit), 5))
	fmt.Printf("primes above 3 up to %d    gaps: %d    decoys: %d\n",
		*limit, len(gaps), *trials)

	fmt.Printf("\n%-4s %-10s %-10s %-10s %-10s %-11s %s\n",
		"k", "free obs", "free r", "div obs", "div r", "PRODUCT", "z(free)")

	rows := []row{}

	for k := 3; k <= *maxK; k += 2 {
		freeObs := pattern.PalindromesWith(gaps, k, centreFree)
		divObs := pattern.PalindromesWith(gaps, k, allDivisible)
		if freeObs == 0 || divObs == 0 {
			fmt.Printf("%-4d %-10d %-10s %-10d %-10s %-11s %s\n",
				k, freeObs, "—", divObs, "—", "—", "empty branch")
			continue
		}

		freeRes := control.Evaluate(freeObs, *trials, func() int {
			return pattern.PalindromesWith(control.ShuffleGaps(gaps, rng), k, centreFree)
		})
		divRes := control.Evaluate(divObs, *trials, func() int {
			return pattern.PalindromesWith(control.ShuffleGaps(gaps, rng), k, allDivisible)
		})

		p := freeRes.Ratio * divRes.Ratio
		rows = append(rows, row{k, freeRes.Ratio, divRes.Ratio, p})

		fmt.Printf("%-4d %-10d %-10.4f %-10d %-10.4f %-11.4f %+.1f\n",
			k, freeObs, freeRes.Ratio, divObs, divRes.Ratio, p, freeRes.ZScore)
	}

	if len(rows) < 2 {
		return
	}

	fmt.Println("\nSTEP FACTORS — is one branch gaining exactly what the other loses?")
	fmt.Printf("%-10s %-12s %-12s %-12s %s\n",
		"k step", "free x", "div x", "free x div", "vs sqrt(2)")
	for i := 1; i < len(rows); i++ {
		a, b := rows[i-1], rows[i]
		fu := b.free / a.free
		fd := b.div / a.div
		fmt.Printf("%-2d -> %-4d %-12.4f %-12.4f %-12.4f %.4f\n",
			a.k, b.k, fu, fd, fu*fd, fu/math.Sqrt2)
	}

	mean, sd := stats(rows)
	fmt.Printf("\nPRODUCT across all k : mean %.4f   sd %.4f   spread %.1f%%\n",
		mean, sd, 100*sd/mean)
	if sd/mean < 0.05 {
		fmt.Println("  -> constant to within 5%. The law moves structure; it does not create it.")
	} else {
		fmt.Println("  -> not constant. The two branches are not reciprocal.")
	}
}

func stats(rows []row) (float64, float64) {
	sum := 0.0
	for _, r := range rows {
		sum += r.product
	}
	mean := sum / float64(len(rows))

	v := 0.0
	for _, r := range rows {
		v += (r.product - mean) * (r.product - mean)
	}
	return mean, math.Sqrt(v / float64(len(rows)))
}
