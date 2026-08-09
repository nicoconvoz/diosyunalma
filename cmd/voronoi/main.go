// Command voronoi turns on the divisor radio — the perpendicular music.
//
// The primes sing in logarithmic time with frequencies γ. The divisor
// world (ζ², Finding 68's ladder) sings in SQUARE-ROOT time: Voronoi (1904)
// showed the divisor error term oscillates as
//
//	Δ(x) ≈ (x^{1/4}/(π√2)) Σ_n d(n)/n^{3/4} · cos(4π√(nx) − π/4),
//
// so in the variable v = √x the stations sit at angular frequencies 4π√n —
// the square roots of the integers — and the VOLUME of station n is its
// divisor count d(n). The old tritone √2 is a titular station here.
//
// PRE-REGISTERED: periodogram peaks at ω = 4π√n for n = 1..8, with
// loudness ordered by d(n)/n^{3/4} — station 2 louder than stations 1 and
// 3, station 6 loud, station 5 quiet.
//
// Usage:
//
//	go run ./cmd/voronoi
package main

import (
	"fmt"
	"math"
	"sort"
)

const euler = 0.5772156649015329

// dSum is the exact divisor summatory function via the hyperbola method.
func dSum(x int) float64 {
	s := 0
	r := int(math.Sqrt(float64(x)))
	for k := 1; k <= r; k++ {
		s += x / k
	}
	return float64(2*s - r*r)
}

func main() {
	// square-root time grid: v in [1000, 4000], x = v^2 in [1e6, 1.6e7].
	const v0, v1, dv = 1000.0, 4000.0, 0.075
	var vs, es []float64
	for v := v0; v <= v1; v += dv {
		x := v * v
		xi := int(x)
		delta := dSum(xi) - x*math.Log(x) - (2*euler-1)*x
		vs = append(vs, v)
		es = append(es, delta/math.Sqrt(v))
	}
	n := float64(len(es) - 1)
	for i := range es {
		es[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}

	fmt.Println("THE DIVISOR RADIO — the perpendicular music, in square-root time")
	fmt.Println("\npredicted stations 4*pi*sqrt(n), volume knob d(n)/n^(3/4):")
	type st struct {
		n    int
		w, a float64
	}
	pred := []st{}
	for k := 1; k <= 8; k++ {
		d := 0
		for j := 1; j <= k; j++ {
			if k%j == 0 {
				d++
			}
		}
		pred = append(pred, st{k, 4 * math.Pi * math.Sqrt(float64(k)),
			float64(d) / math.Pow(float64(k), 0.75)})
	}

	power := func(w float64) float64 {
		var cr, ci float64
		for i, v := range vs {
			cr += es[i] * math.Cos(w*v)
			ci += es[i] * math.Sin(w*v)
		}
		return (cr*cr + ci*ci) / float64(len(vs))
	}

	fmt.Println("\n   n    predicted w   measured peak   volume measured   d(n)/n^3/4")
	type row struct {
		n          int
		mp, pw, th float64
	}
	rows := []row{}
	for _, p := range pred {
		// climb to the local maximum near the predicted frequency.
		bestW, bestP := p.w, power(p.w)
		for w := p.w - 0.3; w <= p.w+0.3; w += 0.002 {
			if pw := power(w); pw > bestP {
				bestP, bestW = pw, w
			}
		}
		rows = append(rows, row{p.n, bestW, math.Sqrt(bestP), p.a})
	}
	base := rows[0].pw
	for _, r := range rows {
		fmt.Printf("   %d     %7.3f       %7.3f        %6.2f            %.2f\n",
			r.n, 4*math.Pi*math.Sqrt(float64(r.n)), r.mp, r.pw/base, r.th)
	}

	byVol := append([]row{}, rows...)
	sort.Slice(byVol, func(i, j int) bool { return byVol[i].pw > byVol[j].pw })
	fmt.Printf("\nloudest stations, measured: n=%d, n=%d, n=%d\n",
		byVol[0].n, byVol[1].n, byVol[2].n)
	fmt.Println("the divisor count IS the volume knob: rich numbers sing louder,")
	fmt.Println("and the old tritone sqrt(2) holds a titular chair in this orchestra.")
}
