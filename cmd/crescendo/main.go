// Command crescendo formalizes the pattern the eye found in the chest.
//
// The reading, in its original phrasing: the residual wave over d = 6..30
// rises softly, dips, barely moves, falls — then the SAME wave repeats over
// d = 36..60 with greater force. If true, the anomaly oscillates with
// period 30 — and 30 = 2·3·5 is the primorial wheel — with amplitude
// growing along d: a crescendo on the wheel's beat.
//
// This command takes the chest's measured 10⁹ table (Finding 54), removes
// the linear background, splits the residual wave into its two period-30
// bars, and measures:
//
//  1. the Pearson correlation between the two bars (self-similarity);
//  2. the gain (amplitude ratio bar2/bar1);
//  3. sign concordance of the period-30 pairs (d, d+30), flagged by
//     whether both members are individually significant;
//  4. the weak extension to the third bar (d = 66..90, sparse data,
//     sign check only).
//
// PRE-REGISTERED for the next telescope run (10^10): the third bar must
// repeat the sign pattern (−, +, −, +, +) with amplitude still growing,
// at full significance — or the crescendo dies.
//
// Usage:
//
//	go run ./cmd/crescendo
package main

import (
	"fmt"
	"math"
)

// the chest's 10^9 measurements (Finding 54): gap, deviation %, sigma %.
type tooth struct {
	d          int
	dev, sigma float64
}

var bar1 = []tooth{
	{6, -0.69, 0.066}, {12, 2.20, 0.146}, {18, 0.14, 0.23},
	{24, -0.10, 0.33}, {30, -2.10, 0.34},
}
var bar2 = []tooth{
	{36, -12.51, 0.64}, {42, -0.03, 0.90}, {48, -9.07, 1.31},
	{54, -5.77, 1.86}, {60, -7.51, 1.90},
}
var bar3 = []tooth{
	{66, -16.48, 3.3}, {72, -10.26, 6.0}, {78, -18.27, 6.8},
	{84, -12.46, 9.6}, {90, 2.65, 13.0},
}

func main() {
	// linear background fitted to the two solid bars.
	all := append(append([]tooth{}, bar1...), bar2...)
	var sx, sy, sxx, sxy float64
	for _, t := range all {
		x := float64(t.d)
		sx += x
		sy += t.dev
		sxx += x * x
		sxy += x * t.dev
	}
	n := float64(len(all))
	slope := (n*sxy - sx*sy) / (n*sxx - sx*sx)
	inter := (sy - slope*sx) / n
	trend := func(d int) float64 { return inter + slope*float64(d) }
	fmt.Printf("THE CRESCENDO — background %.2f %+.4f·d removed\n\n", inter, slope)

	res := func(bar []tooth) []float64 {
		out := make([]float64, len(bar))
		for i, t := range bar {
			out[i] = t.dev - trend(t.d)
		}
		return out
	}
	r1, r2, r3 := res(bar1), res(bar2), res(bar3)

	fmt.Println("  pair (d, d+30)   bar1 resid   bar2 resid   both > 5 sigma?   signs")
	strong, strongOK, weakOK, total := 0, 0, 0, 0
	for i := range r1 {
		s1 := math.Abs(bar1[i].dev-trend(bar1[i].d)) / bar1[i].sigma
		s2 := math.Abs(bar2[i].dev-trend(bar2[i].d)) / bar2[i].sigma
		solid := s1 > 5 && s2 > 5
		agree := (r1[i] > 0) == (r2[i] > 0)
		mark := "discord"
		if agree {
			mark = "CONCORD"
			weakOK++
		}
		total++
		if solid {
			strong++
			if agree {
				strongOK++
			}
		}
		fmt.Printf("  (%2d, %2d)        %+6.2f       %+6.2f       %-5v            %s\n",
			bar1[i].d, bar2[i].d, r1[i], r2[i], solid, mark)
	}

	pearson := func(a, b []float64) float64 {
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
	amp := func(r []float64) float64 {
		s := 0.0
		for _, v := range r {
			s += v * v
		}
		return math.Sqrt(s / float64(len(r)))
	}

	fmt.Printf("\nself-similarity of the two solid bars: correlation %.3f\n", pearson(r1, r2))
	fmt.Printf("gain of the crescendo: amplitude x%.1f (%.2f -> %.2f)\n",
		amp(r2)/amp(r1), amp(r1), amp(r2))
	fmt.Printf("sign concordance: %d/%d overall, %d/%d among fully significant pairs\n",
		weakOK, total, strongOK, strong)

	fmt.Println("\nthird bar (66..90) — sparse, sign check only:")
	fmt.Print("  predicted signs (-, +, -, +, +)   measured residuals: ")
	for _, v := range r3 {
		fmt.Printf(" %+.1f", v)
	}
	fmt.Println()
	fmt.Println("\nPRE-REGISTERED: at 10^10 the third bar must repeat the sign pattern")
	fmt.Println("at full significance with amplitude still growing - or the crescendo dies.")
	fmt.Println("the wave's period is 30 = 2*3*5: the wheel itself keeps the beat.")
}
