// Command oscillation asks whether a hidden wave sits under the ratios.
//
// THE IDEA BEING TESTED. The explicit formula writes every prime count as a
// smooth part plus a sum of cosines in ln x, one per zeta zero. So "a hidden
// sine" is not numerology — it is the actual structure of the primes. The
// question is whether any such wave is visible in THIS observable at THIS
// precision.
//
// PRE-REGISTERED EXPECTATION, stated before the data is looked at. Zeta-zero
// oscillations carry relative amplitude ~1/sqrt(N): about 1e-3 at 10^6 and
// 1e-4 at 10^8. The measurement error per bin here is 1% to 10%. The zeros are
// therefore predicted to be invisible, and this run sets an upper bound. Any
// percent-level oscillation that DID appear would be something else — or noise
// dressed up as a wave.
//
// METHOD. Disjoint quarter-decade bins from 10^6 to 10^8, so the points are
// independent — cumulative prefixes share their windows and smooth any wiggle
// away. Per bin: the centre-free ratio against shuffled decoys. Then a
// two-parameter smooth trend a + b/ln N is fitted (two parameters on eight
// points leaves honest degrees of freedom, unlike a cubic on five), and the
// residuals are examined in units of their own sigma. Chi-square near 1 per
// degree of freedom means the smooth trend suffices: no wave at this
// precision. Alternating signs and chi-square far above 1 would be a wave
// candidate.
//
// Usage:
//
//	go run ./cmd/oscillation [-trials N] [-seed N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/control"
	"github.com/nicoconvoz/diosyunalma/pattern"
	"github.com/nicoconvoz/diosyunalma/primes"
)

func centreFree(w []int) bool { return w[len(w)/2]%3 != 0 }

type bin struct {
	lo, hi  int
	ratio   float64
	sigma   float64 // absolute uncertainty on the ratio
	midLogN float64 // ln of the geometric mid of the bin
}

func main() {
	trials := flag.Int("trials", 30, "decoys per bin")
	seed := flag.Int64("seed", 2026, "random seed")
	flag.Parse()

	if *trials < 5 {
		fmt.Fprintln(os.Stderr, "trials >= 5")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))
	walk := primes.From(primes.Sieve(100_000_000), 5)

	// Quarter-decade edges: 1e6 to 1e8 in factors of 10^0.25.
	edges := []int{}
	for e := 0.0; e <= 8.001; e++ {
		edges = append(edges, int(math.Round(1e6*math.Pow(10, e/4))))
	}

	fmt.Println("PRE-REGISTERED: zeta-zero amplitude ~1/sqrt(N) — 1e-3 to 1e-4 across")
	fmt.Println("these bins, against measurement errors of 1% to 10%. Prediction: not")
	fmt.Println("detectable here. This run sets an upper bound, not a discovery.")

	for _, k := range []int{5, 7} {
		fmt.Printf("\n================ k = %d, disjoint bins ================\n", k)
		fmt.Printf("%-24s %-9s %-11s %-11s %-9s %s\n",
			"bin", "windows", "decoy mean", "ratio", "+/-", "rel err")

		bins := []bin{}
		for i := 0; i+1 < len(edges); i++ {
			lo, hi := edges[i], edges[i+1]
			a := sort.SearchInts(walk, lo)
			b := sort.SearchInts(walk, hi)
			if b-a < 100 {
				continue
			}
			gaps := primes.Gaps(walk[a:b])

			obs := pattern.PalindromesWith(gaps, k, centreFree)
			if obs == 0 {
				fmt.Printf("%-24s %-9d %s\n",
					fmt.Sprintf("[%.1e, %.1e)", float64(lo), float64(hi)), 0, "empty — skipped")
				continue
			}

			res := control.Evaluate(obs, *trials, func() int {
				return pattern.PalindromesWith(control.ShuffleGaps(gaps, rng), k, centreFree)
			})
			if res.Mean == 0 {
				continue
			}

			relObs := 1 / math.Sqrt(float64(obs))
			relDecoy := res.StdDev / (res.Mean * math.Sqrt(float64(*trials)))
			rel := math.Hypot(relObs, relDecoy)

			bins = append(bins, bin{
				lo: lo, hi: hi,
				ratio:   res.Ratio,
				sigma:   res.Ratio * rel,
				midLogN: math.Log(math.Sqrt(float64(lo) * float64(hi))),
			})
			fmt.Printf("%-24s %-9d %-11.1f %-11.4f %-9.4f %.1f%%\n",
				fmt.Sprintf("[%.1e, %.1e)", float64(lo), float64(hi)),
				obs, res.Mean, res.Ratio, res.Ratio*rel, 100*rel)
		}

		if len(bins) < 4 {
			fmt.Println("not enough populated bins to fit a trend")
			continue
		}

		// Weighted least squares for ratio = a + b/lnN.
		var S, Sx, Sy, Sxx, Sxy float64
		for _, p := range bins {
			w := 1 / (p.sigma * p.sigma)
			x := 1 / p.midLogN
			S += w
			Sx += w * x
			Sy += w * p.ratio
			Sxx += w * x * x
			Sxy += w * x * p.ratio
		}
		det := S*Sxx - Sx*Sx
		bCoef := (S*Sxy - Sx*Sy) / det
		aCoef := (Sxx*Sy - Sx*Sxy) / det

		fmt.Printf("\nsmooth trend: ratio = %.4f + %.4f/lnN   (2 params on %d points)\n",
			aCoef, bCoef, len(bins))
		fmt.Printf("%-24s %-11s %-11s %s\n", "bin", "measured", "trend", "residual (sigma)")

		chi2 := 0.0
		signChanges := 0
		prevSign := 0
		maxAbs := 0.0
		for _, p := range bins {
			trend := aCoef + bCoef/p.midLogN
			r := (p.ratio - trend) / p.sigma
			chi2 += r * r
			if math.Abs(r) > maxAbs {
				maxAbs = math.Abs(r)
			}
			sign := 1
			if r < 0 {
				sign = -1
			}
			if prevSign != 0 && sign != prevSign {
				signChanges++
			}
			prevSign = sign
			fmt.Printf("%-24s %-11.4f %-11.4f %+.2f\n",
				fmt.Sprintf("[%.1e, %.1e)", float64(p.lo), float64(p.hi)),
				p.ratio, trend, r)
		}

		dof := len(bins) - 2
		fmt.Printf("\nchi2/dof = %.2f on %d dof    sign changes = %d of %d possible    max |residual| = %.1f sigma\n",
			chi2/float64(dof), dof, signChanges, len(bins)-1, maxAbs)

		switch {
		case chi2/float64(dof) < 2 && maxAbs < 2.5:
			fmt.Println("verdict: the smooth trend suffices. No wave at this precision;")
			fmt.Printf("any hidden oscillation is bounded by ~%.1f%% of the ratio here.\n",
				100*meanSigma(bins))
		case signChanges >= len(bins)-2:
			fmt.Println("verdict: structured alternation beyond the trend — wave candidate,")
			fmt.Println("needs more decoys and an independent seed before it is anything.")
		default:
			fmt.Println("verdict: excess scatter without alternation — unmodelled systematics,")
			fmt.Println("not a wave signature.")
		}
	}
}

func meanSigma(bins []bin) float64 {
	s := 0.0
	for _, p := range bins {
		s += p.sigma / p.ratio
	}
	return s / float64(len(bins))
}
