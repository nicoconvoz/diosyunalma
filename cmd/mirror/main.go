// Command mirror builds the mirror-against-mirror snapshot machine.
//
// The flash, in its original phrasing: a mirror reflecting in another
// mirror, images one inside another until they are lost in the echo — but
// if a snapshot captured the EDGES before the image dissolves, it could
// return specific coordinates far ahead in the record.
//
// That machine exists and is the Riemann–Siegel formula. The infinite sum
// defining Z(t) folds at its self-reflection point N = √(t/2π) — the
// mirror pairing n ↔ t/(2πn) — so infinitely many images collapse into N
// terms, and the dissolving edge is captured by a snapshot correction
// C₀(p). With ~126 terms it computes the stations at height 100,000 —
// two thousand times beyond every radio this laboratory has built.
//
// PRE-REGISTERED: the machine's zeros near t = 100,000 must match the
// published tables on subsequent verification, to a few thousandths.
//
// Usage:
//
//	go run ./cmd/mirror [-height T] [-span S]
package main

import (
	"flag"
	"fmt"
	"math"
)

// theta is the Riemann–Siegel phase.
func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 +
		1/(48*t) + 7/(5760*t*t*t)
}

// z is the Riemann–Siegel Z function: the folded mirror plus the edge
// snapshot (first correction term).
func z(t float64) float64 {
	a := math.Sqrt(t / (2 * math.Pi))
	n := int(a)
	th := theta(t)
	s := 0.0
	for k := 1; k <= n; k++ {
		s += math.Cos(th-t*math.Log(float64(k))) / math.Sqrt(float64(k))
	}
	s *= 2
	// the snapshot of the edges.
	p := a - float64(n)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (n-1)%2 == 1 {
		sign = -1
	}
	return s + sign*math.Pow(t/(2*math.Pi), -0.25)*c0
}

func main() {
	height := flag.Float64("height", 100_000, "look at the record near this height")
	span := flag.Float64("span", 10, "width of the window")
	flag.Parse()

	terms := int(math.Sqrt(*height / (2 * math.Pi)))
	fmt.Printf("THE MIRROR — folded at its self-reflection, %d terms instead of infinity\n", terms)
	fmt.Printf("\nstations found near height %.0f:\n", *height)

	prevT := *height
	prevZ := z(prevT)
	count := 0
	for t := *height + 0.005; t <= *height+*span; t += 0.005 {
		zt := z(t)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for hi-lo > 1e-9 {
				mid := (lo + hi) / 2
				if (z(lo) < 0) != (z(mid) < 0) {
					hi = mid
				} else {
					lo = mid
				}
			}
			count++
			fmt.Printf("   %.6f\n", (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	fmt.Printf("\n%d stations in a window of %.0f, from a folded sum of %d terms.\n",
		count, *span, terms)
	fmt.Println("the images cancel in pairs at the mirror point; the snapshot C0 keeps")
	fmt.Println("the edges. specific coordinates, far ahead in the record - the flash, built.")
}
