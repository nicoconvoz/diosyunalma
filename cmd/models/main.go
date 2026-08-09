// Command models decides between functional forms by prediction, not by fit.
//
// The centre-free ratio has been measured at k = 3, 5, 7, 9 and 11. Several
// shapes can be laid over those points, and the ones with more parameters will
// always look better:
//
//	geometric   2 parameters   determined by k = 3, 5
//	linear      2 parameters   determined by k = 3, 5
//	quadratic   3 parameters   determined by k = 3, 5, 7
//	cubic       4 parameters   determined by k = 3, 5, 7, 9
//
// A cubic through four points fits them exactly whatever they are, including
// pure noise. Its residual is zero by construction and carries no information.
//
// So each model is given exactly the points that determine it, and then made to
// predict k = 11, which none of them has seen. The measured value at k = 11
// carries its own error bar, and a model is only wrong if it misses by more
// than that.
//
// Usage:
//
//	go run ./cmd/models [-limit N] [-trials N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/nicoconvoz/numerosprimos/control"
	"github.com/nicoconvoz/numerosprimos/pattern"
	"github.com/nicoconvoz/numerosprimos/primes"
)

func centreFree(w []int) bool { return w[len(w)/2]%3 != 0 }

type point struct {
	k        int
	observed int
	ratio    float64
	relErr   float64 // fractional uncertainty on the ratio
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	trials := flag.Int("trials", 60, "decoys per window size")
	seed := flag.Int64("seed", 2026, "random seed")
	flag.Parse()

	if *limit < 1000 || *trials < 5 {
		fmt.Fprintln(os.Stderr, "limit >= 1000 and trials >= 5")
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(*seed))
	gaps := primes.Gaps(primes.From(primes.Sieve(*limit), 5))
	fmt.Printf("primes above 3 up to %d    gaps: %d    decoys: %d\n\n",
		*limit, len(gaps), *trials)

	fmt.Println("MEASURED — centre-free ratio with its uncertainty")
	fmt.Printf("%-4s %-11s %-12s %-11s %s\n", "k", "observed", "decoy mean", "ratio", "+/-")

	pts := []point{}
	for k := 3; k <= 11; k += 2 {
		obs := pattern.PalindromesWith(gaps, k, centreFree)
		if obs == 0 {
			continue
		}
		res := control.Evaluate(obs, *trials, func() int {
			return pattern.PalindromesWith(control.ShuffleGaps(gaps, rng), k, centreFree)
		})
		if res.Mean == 0 {
			continue
		}

		// Counting error on the observation, plus the error on the decoy mean.
		relObs := 1 / math.Sqrt(float64(obs))
		relDecoy := 0.0
		if res.StdDev > 0 {
			relDecoy = res.StdDev / (res.Mean * math.Sqrt(float64(*trials)))
		}
		rel := math.Hypot(relObs, relDecoy)

		pts = append(pts, point{k, obs, res.Ratio, rel})
		fmt.Printf("%-4d %-11d %-12.1f %-11.4f %.2f%%\n",
			k, obs, res.Mean, res.Ratio, 100*rel)
	}

	if len(pts) < 5 {
		fmt.Fprintln(os.Stderr, "\nneed all five window sizes to run the comparison")
		os.Exit(1)
	}

	target := pts[len(pts)-1] // k = 11, held out from every model

	fmt.Printf("\nHELD OUT: k=%d measured %.4f +/- %.4f (%.1f%%)\n",
		target.k, target.ratio, target.ratio*target.relErr, 100*target.relErr)

	fmt.Println("\nPREDICTIONS — each model sees only the points that determine it")
	fmt.Printf("%-12s %-7s %-13s %-12s %-11s %s\n",
		"model", "params", "fitted on", "predicts k=11", "error", "verdict")

	fit := pts[:4] // k = 3, 5, 7, 9

	report := func(name string, params int, on string, pred float64) {
		err := (pred - target.ratio) / target.ratio
		sigma := math.Abs(err) / target.relErr
		verdict := "survives"
		if sigma > 2 {
			verdict = fmt.Sprintf("REJECTED at %.1f sigma", sigma)
		}
		fmt.Printf("%-12s %-7d %-13s %-12.4f %-11s %s\n",
			name, params, on, pred, fmt.Sprintf("%+.1f%%", 100*err), verdict)
	}

	// Geometric: constant factor per step, anchored on the first two points.
	step := fit[1].ratio / fit[0].ratio
	report("geometric", 2, "k=3,5", fit[0].ratio*math.Pow(step, float64(target.k-3)/2))

	// Linear in k through the first two points.
	slope := (fit[1].ratio - fit[0].ratio) / 2
	report("linear", 2, "k=3,5", fit[0].ratio+slope*float64(target.k-3))

	// Quadratic through the first three, cubic through all four.
	report("quadratic", 3, "k=3,5,7", lagrange(fit[:3], float64(target.k)))
	report("cubic", 4, "k=3,5,7,9", lagrange(fit[:4], float64(target.k)))

	// How far each fitted point sits from the geometric law, in its own sigma.
	fmt.Println("\nRESIDUALS AGAINST THE GEOMETRIC LAW — is any deviation real?")
	fmt.Printf("%-4s %-11s %-11s %-11s %s\n", "k", "measured", "geometric", "deviation", "sigma")
	for _, p := range pts {
		pred := fit[0].ratio * math.Pow(step, float64(p.k-3)/2)
		dev := (p.ratio - pred) / pred
		fmt.Printf("%-4d %-11.4f %-11.4f %-11s %.1f\n",
			p.k, p.ratio, pred, fmt.Sprintf("%+.1f%%", 100*dev),
			math.Abs(dev)/p.relErr)
	}
	fmt.Println("\nA deviation under 2 sigma is not evidence of a different shape.")
}

// lagrange interpolates the polynomial through every given point and evaluates
// it at x. With n points it is the unique degree n-1 polynomial through them,
// so its residual on those points is zero whatever they contain.
func lagrange(pts []point, x float64) float64 {
	sum := 0.0
	for i, pi := range pts {
		term := pi.ratio
		for j, pj := range pts {
			if i == j {
				continue
			}
			term *= (x - float64(pj.k)) / float64(pi.k-pj.k)
		}
		sum += term
	}
	return sum
}
