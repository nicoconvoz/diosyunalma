// Command rhythm asks whether the orchestra keeps one beat.
//
// The laboratory's own measured zeros — nine dials, some sixty stations —
// allow two structural questions, both pre-registered:
//
//  1. WITHIN each dial: pooled unfolded spacings should REPEL (few small
//     gaps, variance far below 1) — the universality conjecture says every
//     L-function drums to the same GUE rhythm the zeta dial showed in
//     Finding 28.
//
//  2. ACROSS dials: zeros of distinct primitive L-functions are conjectured
//     statistically independent — merged lists should show NO repulsion,
//     with close encounters between different tribes' stations occurring
//     freely, pushing the merged statistics toward Poisson.
//
// Every musician keeps strict time with itself and ignores the player beside
// it — which is why the orchestra sounds rich instead of military.
package main

import (
	"fmt"
	"math"
	"sort"
)

type dial struct {
	name string
	q    float64 // conductor, entering the zero-density law
	g    []float64
}

// our own measured stations (Findings 26, 41-44).
var orchestra = []dial{
	{"zeta", 1, []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422, 37.5872, 40.9264, 43.3211, 48.0105, 49.7752}},
	{"chi3", 3, []float64{8.0396, 11.2450, 15.7062, 18.2579, 20.4551, 24.0636}},
	{"chi4", 4, []float64{6.0199, 10.2423, 12.9848, 16.3464, 18.2914, 21.4547}},
	{"chi5", 5, []float64{6.6516, 9.8280, 11.9612, 16.0386, 17.5632, 19.5431, 22.2228, 24.5864}},
	{"chi7", 7, []float64{4.4762, 6.8400, 11.1782, 12.4670, 15.1161, 16.7892}},
	{"chi8", 8, []float64{4.8989, 7.6194, 10.8219, 12.3126, 15.1935, 17.0246}},
	{"chi11", 11, []float64{2.4768, 6.7997, 8.9663, 10.1257, 13.0422, 15.0996}},
	{"chi13", 13, []float64{3.1119, 7.2340, 8.6013, 10.3241, 12.6185, 15.1341}},
	{"chi5 complex +", 5, []float64{6.1829, 8.4594, 12.6747, 14.8310}},
	{"chi5 complex -", 5, []float64{4.1322, 9.4433, 11.2861, 14.1115}},
}

// density is the zero-counting density of an L-function of conductor q at
// height t: ln(q·t/2π)/2π per unit height.
func density(q, t float64) float64 {
	return math.Log(q*t/(2*math.Pi)) / (2 * math.Pi)
}

func stats(sp []float64) (mean, variance, small float64) {
	for _, s := range sp {
		mean += s
	}
	mean /= float64(len(sp))
	n := 0
	for _, s := range sp {
		variance += (s - mean) * (s - mean)
		if s < 0.5*mean {
			n++
		}
	}
	variance /= float64(len(sp)) * mean * mean // variance of s/mean
	small = float64(n) / float64(len(sp))
	return
}

func main() {
	fmt.Println("THE RHYTHM OF THE ORCHESTRA — our own sixty stations")

	// 1) WITHIN: pooled unfolded consecutive spacings per dial.
	within := []float64{}
	for _, d := range orchestra {
		for i := 1; i < len(d.g); i++ {
			mid := (d.g[i] + d.g[i-1]) / 2
			if mid < 2 {
				continue
			}
			within = append(within, (d.g[i]-d.g[i-1])*density(d.q, mid))
		}
	}
	m, v, small := stats(within)
	fmt.Printf("\n1) WITHIN each dial — %d spacings pooled across nine instruments\n", len(within))
	fmt.Printf("   mean %.3f   var %.3f   frac < half-mean %.1f%%\n", m, v, 100*small)
	fmt.Printf("   GUE rhythm:      var 0.178   small-gap  11%%\n")
	fmt.Printf("   Poisson (no rhythm): var 1.000   small-gap  39%%\n")

	// 2) ACROSS: merge every dial in a common window and measure the merged
	// spacings. Independence pushes a superposition toward Poisson.
	type station struct{ g, dens float64 }
	merged := []station{}
	for _, d := range orchestra {
		for _, g := range d.g {
			if g >= 4 && g <= 25 {
				merged = append(merged, station{g, density(d.q, g)})
			}
		}
	}
	sort.Slice(merged, func(i, j int) bool { return merged[i].g < merged[j].g })

	// total density of the superposed process = sum of the dials' densities.
	totalDens := func(t float64) float64 {
		s := 0.0
		for _, d := range orchestra {
			s += density(d.q, t)
		}
		return s
	}
	across := []float64{}
	closest, cA, cB := math.Inf(1), "", ""
	for i := 1; i < len(merged); i++ {
		mid := (merged[i].g + merged[i-1].g) / 2
		gap := (merged[i].g - merged[i-1].g) * totalDens(mid)
		across = append(across, gap)
		if raw := merged[i].g - merged[i-1].g; raw < closest {
			closest = raw
			cA = fmt.Sprintf("%.4f", merged[i-1].g)
			cB = fmt.Sprintf("%.4f", merged[i].g)
		}
	}
	m2, v2, small2 := stats(across)
	fmt.Printf("\n2) ACROSS dials — %d spacings of the merged orchestra in [4, 25]\n", len(across))
	fmt.Printf("   mean %.3f   var %.3f   frac < half-mean %.1f%%\n", m2, v2, 100*small2)
	fmt.Printf("   closest encounter between different tribes: %s and %s (raw gap %.4f)\n",
		cA, cB, closest)

	fmt.Println("\nPRE-REGISTERED READINGS")
	fmt.Println("  within repels (var well below 1, few small gaps) -> one shared rhythm")
	fmt.Println("  across relaxes toward Poisson -> the musicians ignore each other")
}
