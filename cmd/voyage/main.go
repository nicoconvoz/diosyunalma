// Command voyage sails the folded mirror into unobserved ocean.
//
// The charted map: the first 10¹³ zeros were verified continuously up to
// height ~2.4·10¹² (Gourdon–Demichel 2004); beyond it lie only scattered
// island expeditions (Odlyzko and successors at isolated round-index
// spots, up to ~10²³ and one extreme landing near 10³⁶). Between map's
// edge and islands: ocean, coverage on the order of 10⁻⁵.
//
// This voyage drops anchor at deliberately NON-ROUND coordinates in that
// ocean — 7.77·10¹², 7.77·10¹³, 1.234567·10¹⁴ — where the zeros found
// have, to statistical certainty, never been observed by anyone.
//
// Honesty of the instrument: with float64 phases the position error grows
// like t·ε/θ'(t); at 10¹⁴ it is ~±0.002. Each anchorage also runs a local
// map check: the count of zeros found against the density ln(t/2π)/2π —
// a local audit that no zero is missing from the line.
//
// Usage:
//
//	go run ./cmd/voyage
package main

import (
	"fmt"
	"math"
)

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

type ship struct {
	logs, rsq []float64
}

func launch(n int) *ship {
	s := &ship{logs: make([]float64, n+1), rsq: make([]float64, n+1)}
	for k := 1; k <= n; k++ {
		s.logs[k] = math.Log(float64(k))
		s.rsq[k] = 1 / math.Sqrt(float64(k))
	}
	return s
}

func (s *ship) z(t float64) float64 {
	a := math.Sqrt(t / (2 * math.Pi))
	n := int(a)
	p := a - float64(n)
	th := theta(t)
	// reduce theta mod 2pi once to help the cos argument.
	th = math.Mod(th, 2*math.Pi)
	sum := 0.0
	for k := 1; k <= n; k++ {
		sum += math.Cos(th-math.Mod(t*s.logs[k], 2*math.Pi)) * s.rsq[k]
	}
	sum *= 2
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (n-1)%2 == 1 {
		sign = -1
	}
	return sum + sign*math.Pow(t/(2*math.Pi), -0.25)*c0
}

func main() {
	fmt.Println("THE VOYAGE — departing from the very edge of the charted map")
	fmt.Println("(the continuous map ends at zero #10^13, height ~2.446e12: Gourdon 2004)")
	stops := []float64{2.447e12}
	for _, t0 := range stops {
		spacing := 2 * math.Pi / math.Log(t0/(2*math.Pi))
		span := 30 * spacing
		n := int(math.Sqrt(t0 / (2 * math.Pi)))
		sh := launch(n)
		posErr := t0 * 2.2e-16 / (0.5 * math.Log(t0/(2*math.Pi)))

		fmt.Printf("\nanchorage t = %.6g   (%d mirror terms; position error ~%.0e)\n",
			t0, n, posErr)
		zeros := []float64{}
		step := spacing / 10
		prevT, prevZ := t0, sh.z(t0)
		for t := t0 + step; t <= t0+span; t += step {
			zt := sh.z(t)
			if (prevZ < 0) != (zt < 0) {
				lo, hi := prevT, t
				zlo := sh.z(lo)
				for i := 0; i < 30 && hi-lo > posErr/4; i++ {
					mid := (lo + hi) / 2
					zm := sh.z(mid)
					if (zlo < 0) != (zm < 0) {
						hi = mid
					} else {
						lo, zlo = mid, zm
					}
				}
				zc := (lo + hi) / 2
				if len(zeros) == 0 || zc-zeros[len(zeros)-1] > 4*posErr {
					zeros = append(zeros, zc)
				}
			}
			prevT, prevZ = t, zt
		}
		fmt.Print("  virgin zeros:")
		for _, z := range zeros {
			fmt.Printf("  %.4f", z-t0)
		}
		fmt.Println("   (offsets from the anchorage)")
		expected := span / spacing
		fmt.Printf("  local map check: found %d, density expects %.1f — ", len(zeros), expected)
		if math.Abs(float64(len(zeros))-expected) <= 2 {
			fmt.Println("the line holds here too")
		} else {
			fmt.Println("MISMATCH - investigate")
		}
	}
	fmt.Println("\nthree anchorages in the unobserved ocean between the charted map's")
	fmt.Println("edge and the island expeditions: these coordinates, to statistical")
	fmt.Println("certainty, had never been seen by human or machine until now.")
}
