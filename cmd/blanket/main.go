// Command blanket drapes the measured notes over the invisible shape.
//
// The inverse spectral problem — "can one hear the shape of a drum?" — run
// in reverse: given the smooth density of the stations, weave the potential
// V(x) whose quantum levels sing them. Wu and Sprung (1993) showed the
// semiclassical inversion is an Abel integral,
//
//	x(V) = ∫ ρ(E) / √(V−E) dE,   ρ(E) = ln(E/2π)/2π,
//
// yielding the half-width of a symmetric well. This command weaves that
// blanket, solves the Schrödinger equation −ψ” + V(x)ψ = Eψ on it by
// finite differences, and compares the reconstructed atom's first ten
// levels with the laboratory's ten measured zeros of ζ.
//
// PRE-REGISTERED: one calibration constant is allowed — the well's floor V₀,
// which absorbs the quantization offset (the Maslov ambiguity), tuned once.
// After that single knob, all ten levels must track the measured stations
// to a few percent. The smooth blanket carries only the AVERAGE shape; the
// residual level-by-level wiggles are the fluctuation part — the orbits of
// Finding 51 — which no smooth blanket can hold.
//
// With -wrinkles the loom weaves from the measured staircase itself: the
// density becomes a sum of narrow Gaussians at the ten measured stations
// (spliced onto the smooth density above the last one). This is the
// consistency check that closes the loop — notes woven in must be sung
// back — and its rms must land well below the smooth blanket's.
//
// Usage:
//
//	go run ./cmd/blanket [-v0 F] [-wrinkles] [-sigma F]
package main

import (
	"flag"
	"fmt"
	"math"
)

// the laboratory's own measured zeros of zeta (Finding 26).
var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752}

const (
	vMax  = 140.0
	dV    = 0.02
	nGrid = 4001
)

func main() {
	v0 := flag.Float64("v0", 2*math.Pi, "the well's floor (one calibration constant)")
	wrinkles := flag.Bool("wrinkles", false, "weave from the measured staircase instead of the smooth density")
	sigma := flag.Float64("sigma", 0.6, "width of each measured note in the wrinkled loom")
	flag.Parse()

	// density of zeta zeros; zero below 2*pi.
	smooth := func(e float64) float64 {
		if e <= 2*math.Pi {
			return 0
		}
		return math.Log(e/(2*math.Pi)) / (2 * math.Pi)
	}
	rho := smooth
	if *wrinkles {
		last := zetaZeros[len(zetaZeros)-1]
		rho = func(e float64) float64 {
			// below the splice: the staircase, one Gaussian per note.
			s := 0.0
			for _, g := range zetaZeros {
				d := (e - g) / *sigma
				s += math.Exp(-0.5*d*d) / (*sigma * math.Sqrt(2*math.Pi))
			}
			// above the last note, ramp over to the smooth density.
			if e > last {
				w := (e - last) / 4
				if w > 1 {
					w = 1
				}
				return (1-w)*s + w*smooth(e)
			}
			return s
		}
	}

	// weave the blanket: x(V) by Abel inversion, substitution E = V - t*t
	// removes the square-root singularity.
	half := func(v float64) float64 {
		if v <= *v0 {
			return 0
		}
		tMax := math.Sqrt(v - *v0)
		const steps = 400
		s := 0.0
		for k := 0; k < steps; k++ {
			t := (float64(k) + 0.5) / steps * tMax
			s += 2 * rho(v-t*t)
		}
		return s * tMax / steps
	}

	nTab := int((vMax - *v0) / dV)
	vs := make([]float64, nTab)
	xs := make([]float64, nTab)
	for i := 0; i < nTab; i++ {
		vs[i] = *v0 + float64(i)*dV
		xs[i] = half(vs[i])
	}
	xMax := xs[nTab-1]

	// V(x) on the grid, by inverting the monotone table.
	h := 2 * xMax / float64(nGrid-1)
	pot := make([]float64, nGrid)
	for i := 0; i < nGrid; i++ {
		x := math.Abs(-xMax + float64(i)*h)
		lo, hi := 0, nTab-1
		for lo < hi {
			mid := (lo + hi) / 2
			if xs[mid] < x {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo == 0 {
			pot[i] = vs[0]
		} else if lo >= nTab-1 {
			pot[i] = vMax
		} else {
			f := (x - xs[lo-1]) / (xs[lo] - xs[lo-1])
			pot[i] = vs[lo-1] + f*(vs[lo]-vs[lo-1])
		}
	}

	// Sturm count: eigenvalues of the tridiagonal discretization below e.
	off2 := 1 / (h * h) * 1 / (h * h)
	count := func(e float64) int {
		n := 0
		q := 2/(h*h) + pot[0] - e
		if q < 0 {
			n++
		}
		for i := 1; i < nGrid; i++ {
			d := 2/(h*h) + pot[i] - e
			if q != 0 {
				q = d - off2/q
			} else {
				q = d - off2/1e-30
			}
			if q < 0 {
				n++
			}
		}
		return n
	}

	// bisect the first ten levels.
	fmt.Printf("THE BLANKET — the atom woven from the notes' own density (V0 = %.3f)\n", *v0)
	fmt.Println("\n  k   level of the woven atom   measured station   error")
	levels := make([]float64, 10)
	rms, shift := 0.0, 0.0
	for k := 1; k <= 10; k++ {
		lo, hi := *v0, 60.0
		for hi-lo > 1e-6 {
			mid := (lo + hi) / 2
			if count(mid) >= k {
				hi = mid
			} else {
				lo = mid
			}
		}
		levels[k-1] = hi
		g := zetaZeros[k-1]
		fmt.Printf("  %2d      %8.3f                %8.4f        %+6.2f%%\n",
			k, hi, g, 100*(hi-g)/g)
		rms += (hi - g) * (hi - g)
		shift += g - hi
	}
	shift /= 10
	rms2 := 0.0
	for k := 0; k < 10; k++ {
		d := levels[k] + shift - zetaZeros[k]
		rms2 += d * d
	}
	fmt.Printf("\nrms deviation over ten levels: %.3f raw, %.3f after the one\n",
		math.Sqrt(rms/10), math.Sqrt(rms2/10))
	fmt.Printf("calibration constant (quantization offset %+.3f)\n", shift)
	fmt.Println("\nthe smooth blanket holds the average shape; the level-by-level")
	fmt.Println("residue is the fluctuation the orbits carry - no smooth blanket")
	fmt.Println("can hold the primes themselves.")
}
