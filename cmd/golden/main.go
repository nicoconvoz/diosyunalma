// Command golden scans every constant this laboratory has measured or derived
// against the golden-ratio family, with the honest yardstick built in.
//
// The trap, quantified before looking: the family {1/φ², 1/φ, φ/2, 2/φ, √φ,
// 2−1/φ, φ, φ²} carpets the working range so densely that a RANDOM number
// lands within a few percent of some member. The test is therefore not "is
// anything close" — something always is — but "is our set closer than chance,
// and are the exactly-known constants exactly golden".
package main

import (
	"fmt"
	"math"
	"math/rand"
)

var phi = (1 + math.Sqrt(5)) / 2

type target struct {
	name string
	v    float64
}

type constant struct {
	name  string
	v     float64
	exact bool // known to many digits (derived), so near-miss = definitive miss
}

func main() {
	family := []target{
		{"1/phi^2", 1 / (phi * phi)},
		{"1/phi", 1 / phi},
		{"phi/2", phi / 2},
		{"2/phi", 2 / phi},
		{"sqrt(phi)", math.Sqrt(phi)},
		{"2-1/phi", 2 - 1/phi},
		{"phi", phi},
		{"phi^2", phi * phi},
	}

	consts := []constant{
		{"Euler product C (F20)", 0.81980245, true},
		{"c0/c destination (F35)", 2.0 / 3.0, true},
		{"step destination 4/3 (F35)", 4.0 / 3.0, true},
		{"mean beta (F27)", 0.4995, false},
		{"P(stay) (F13)", 0.43036, false},
		{"c0/c at 1e9 (F35)", 0.5638, false},
		{"step factor at 1e8 (F22)", 1.4044, false},
		{"odd ratio at 1e8 (F31)", 1.224, false},
		{"even ratio k=2 (F31)", 0.456, false},
		{"free-branch ratio k=3 (F21)", 1.4443, false},
		{"div-branch ratio k=3 (F21)", 0.8511, false},
		{"outer-pair class 0 (A1)", 0.4009, false},
		{"GUE variance (F28)", 0.178, false},
		{"survival exponent (F18)", 0.166, false},
		{"bits/prime gaps (F19)", 4.1708 / 3, false}, // scaled into range
	}

	fmt.Println("THE SCAN — each constant vs its nearest golden relative")
	fmt.Printf("%-30s %-10s %-12s %-10s %s\n",
		"constant", "value", "nearest", "distance", "verdict")

	sum := 0.0
	for _, c := range consts {
		bestName, bestDist := "", math.Inf(1)
		for _, t := range family {
			d := math.Abs(c.v-t.v) / t.v
			if d < bestDist {
				bestDist, bestName = d, t.name
			}
		}
		sum += bestDist
		verdict := "carpet noise"
		if bestDist < 0.005 {
			verdict = "CLOSE - examine"
		}
		if c.exact && bestDist > 0.005 {
			verdict = "exact constant: definitively NOT golden"
		}
		fmt.Printf("%-30s %-10.4f %-12s %-10.2f%% %s\n",
			c.name, c.v, bestName, 100*bestDist, verdict)
	}
	fmt.Printf("\nmean distance of OUR constants : %.2f%%\n", 100*sum/float64(len(consts)))

	// The carpet density: what a random number scores against the same family.
	rng := rand.New(rand.NewSource(1))
	trials := 200000
	rsum := 0.0
	for i := 0; i < trials; i++ {
		v := 0.15 + rng.Float64()*2.5
		best := math.Inf(1)
		for _, t := range family {
			if d := math.Abs(v-t.v) / t.v; d < best {
				best = d
			}
		}
		rsum += best
	}
	fmt.Printf("mean distance of RANDOM numbers: %.2f%%   <- the carpet\n",
		100*rsum/float64(trials))

	// The bus route: stations the drifting c0/c crosses on its way to 2/3.
	fmt.Println("\nTHE BUS ROUTE of c0/c (measured 0.527 -> 0.564, destination 2/3):")
	fmt.Printf("  tritone station  2-sqrt2 = %.5f   crossed in transit\n", 2-math.Sqrt2)
	fmt.Printf("  GOLDEN station   1/phi   = %.5f   ALSO on the route, after the tritone\n", 1/phi)
	fmt.Printf("  terminal         2/3     = %.5f\n", 2.0/3.0)
}
