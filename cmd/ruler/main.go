// Command ruler formalizes the reading of the three-decade landscape.
//
// The observations, in their original phrasing: the three curves have
// almost the same shape up to 48; each is a pulsar of a different length;
// the lines in panel two move with lives of their own; and the blue one —
// gap 12, frozen — must be THE RULER, the most perfect of all, the one to
// measure with.
//
// Formalized, that is a SCALING HYPOTHESIS: the deviation profile has one
// fixed shape across decades of x, with only its amplitude fading, anchored
// by the scale-free 12. This command measures, on the recorded tables of
// Findings 54 and 63:
//
//  1. shape correlation between decades (window 6..60);
//  2. amplitude decay per decade;
//  3. each tooth's private decay speed (the "life of its own");
//  4. the ruler: every deviation re-measured in units of the invariant.
//
// PRE-REGISTERED for 10^11: shape correlation with the 10^10 profile must
// stay above 0.9, the amplitude must fade by roughly x0.8 again, and the
// ruler itself must hold at +2.2% — or the scaling law dies.
//
// Usage:
//
//	go run ./cmd/ruler
package main

import (
	"fmt"
	"math"
)

var gaps = []int{6, 12, 18, 24, 30, 36, 42, 48, 54, 60}

var dev = map[int][]float64{
	8:  {-0.97, 1.92, -0.70, -0.47, -6.08, -12.13, 1.92, -7.12, -18.92, -9.24},
	9:  {-0.69, 2.20, 0.14, -0.10, -2.10, -12.51, -0.03, -9.07, -5.77, -7.51},
	10: {-0.53, 2.21, 0.47, 0.11, -1.11, -10.72, -0.37, -6.03, -2.34, -4.86},
}

func pearson(a, b []float64) float64 {
	var ma, mb float64
	for i := range a {
		ma += a[i]
		mb += b[i]
	}
	ma /= float64(len(a))
	mb /= float64(len(b))
	var num, da, db float64
	for i := range a {
		num += (a[i] - ma) * (b[i] - mb)
		da += (a[i] - ma) * (a[i] - ma)
		db += (b[i] - mb) * (b[i] - mb)
	}
	return num / math.Sqrt(da*db)
}

func amp(a []float64) float64 {
	var m float64
	for _, v := range a {
		m += v
	}
	m /= float64(len(a))
	var s float64
	for _, v := range a {
		s += (v - m) * (v - m)
	}
	return math.Sqrt(s / float64(len(a)))
}

func main() {
	fmt.Println("THE RULER — the scaling law behind the three-decade landscape")

	fmt.Println("\n1) one shape, fading volume (window 6..60):")
	fmt.Printf("   shape correlation 10^8 <-> 10^9:  %.2f\n", pearson(dev[8], dev[9]))
	fmt.Printf("   shape correlation 10^9 <-> 10^10: %.2f\n", pearson(dev[9], dev[10]))
	fmt.Printf("   amplitude per decade: %.2f -> %.2f -> %.2f  (fade x%.2f, x%.2f)\n",
		amp(dev[8]), amp(dev[9]), amp(dev[10]),
		amp(dev[9])/amp(dev[8]), amp(dev[10])/amp(dev[9]))

	fmt.Println("\n2) each tooth's private decay (dev ratio per decade) — the life of its own:")
	for i, d := range gaps {
		if dev[10][i] > -1 && d != 36 {
			continue
		}
		r1 := dev[9][i] / dev[8][i]
		r2 := dev[10][i] / dev[9][i]
		fmt.Printf("   gap %2d: x%.2f then x%.2f\n", d, r1, r2)
	}

	fmt.Println("\n3) the ruler — every 10^10 deviation in units of the invariant (+2.21%):")
	for i, d := range gaps {
		fmt.Printf("   gap %2d: %+6.2f rulers\n", d, dev[10][i]/2.21)
	}

	fmt.Println("\nPRE-REGISTERED for 10^11: shape correlation vs 10^10 above 0.9,")
	fmt.Println("amplitude fading near x0.8, the ruler still +2.2 - or the law dies.")
	fmt.Println("one melody, fading volume, one frozen note to measure it all by.")
}
