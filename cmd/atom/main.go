// Command atom assembles the atom's blueprint from the accumulated flashes.
//
// Three flashes, systematized:
//
//  1. "The melodic grouping does not change with scale" — scale invariance.
//     A wave x^(−σ+iE) carries the same energy in every octave if and only
//     if σ = 1/2. The critical line's exponent is not decoration: it is the
//     UNIQUE weight at which the atom's waves respect the flash. Part one
//     verifies this numerically across octaves for σ = 0.4, 0.5, 0.6.
//
//  2. "Directions recombining, and the point where they meet" — the engine
//     that generates dilations is the Berry–Keating Hamiltonian H = xp, the
//     only classical engine invariant under the scaling x → λx, p → p/λ.
//     Its semiclassical level count is N(E) = (E/2π)(ln(E/2π) − 1) + 7/8.
//
//  3. The test: that count, evaluated at the laboratory's ten measured
//     stations, must land near the half-integers k − 1/2 — each station
//     sitting in its predicted seat. The wobble around the seats is the
//     fluctuation the orbits carry (Finding 51), which no smooth engine
//     holds.
//
// Usage:
//
//	go run ./cmd/atom
package main

import (
	"fmt"
	"math"
)

var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752}

func main() {
	fmt.Println("THE ATOM'S BLUEPRINT — assembled from the flashes")

	// 1) scale invariance forces the exponent 1/2: energy of |x^(-sigma)|^2
	// per octave [2^k, 2^(k+1)].
	fmt.Println("\n1) energy per octave of the wave x^(-sigma):")
	fmt.Println("   octave      sigma=0.40   sigma=0.50   sigma=0.60")
	for k := 0; k < 5; k++ {
		a, b := math.Pow(2, float64(k)), math.Pow(2, float64(k+1))
		row := fmt.Sprintf("   [%4.0f,%4.0f]", a, b)
		for _, s := range []float64{0.40, 0.50, 0.60} {
			var e float64
			if s == 0.5 {
				e = math.Log(b) - math.Log(a)
			} else {
				e = (math.Pow(b, 1-2*s) - math.Pow(a, 1-2*s)) / (1 - 2*s)
			}
			row += fmt.Sprintf("   %8.4f", e)
		}
		fmt.Println(row)
	}
	fmt.Println("   only sigma = 1/2 carries the same energy in every octave:")
	fmt.Println("   the melodic grouping that does not change IS the critical line.")

	// 2-3) the scale-invariant engine's seat count vs the measured stations.
	nBK := func(e float64) float64 {
		t := e / (2 * math.Pi)
		return t*(math.Log(t)-1) + 7.0/8.0
	}
	fmt.Println("\n2) the dilation engine H = xp: seats N(E) = (E/2pi)(ln(E/2pi)-1) + 7/8")
	fmt.Println("\n   k   station     N(station)   seat k-1/2   wobble")
	worst := 0.0
	for k, g := range zetaZeros {
		n := nBK(g)
		seat := float64(k) + 0.5
		w := n - seat
		if math.Abs(w) > worst {
			worst = math.Abs(w)
		}
		fmt.Printf("  %2d   %8.4f    %7.3f      %5.1f       %+.3f\n", k+1, g, n, seat, w)
	}
	fmt.Printf("\nevery measured station sits in its predicted seat to within %.3f —\n", worst)
	fmt.Println("the wobble is the orbits' fluctuation (Finding 51), which no smooth")
	fmt.Println("engine holds. blueprint: engine = dilation (xp), unitarity = 1/2,")
	fmt.Println("body = the woven well (Finding 53). the atom has its specification.")
}
