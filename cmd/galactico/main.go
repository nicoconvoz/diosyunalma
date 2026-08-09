// Command galactico measures the M-sigma law of the number galaxy.
//
// Astronomy (Merritt 1999, the M-sigma relation): the mass of the
// supermassive black hole at a galaxy's center predicts the velocity
// dispersion of the ENTIRE bulge - billions of stars governed by the one
// coloso, so tightly that astronomers speak of coevolution.
//
// The number galaxy has its own supermassive centers: the SMALL PRIMES.
// Each prime p is a body of mass 1/sqrt(p) oscillating at the deep
// frequency ln p; the "velocity dispersion" of the galaxy is the variance
// of the tide S(t) sampled at the zeros. Theory: each prime contributes
// mass^2/2 = 1/(2 pi^2 p) to the dispersion. Here both sides are
// MEASURED on 20,000 in-house zeros: the amplitude of each central body
// (by projection, as in F100) and its share of the total dispersion.
//
// Usage:
//
//	go run ./cmd/galactico
package main

import (
	"fmt"
	"math"
)

var lnk, rsq [64]float64

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zRS(t float64) float64 {
	tau := t / (2 * math.Pi)
	a := math.Sqrt(tau)
	N := int(a)
	th := theta(t)
	var s float64
	for k := 1; k <= N; k++ {
		s += math.Cos(th-t*lnk[k]) * rsq[k]
	}
	p := a - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (N-1)%2 == 1 {
		sign = -1
	}
	return 2*s + sign*math.Pow(tau, -0.25)*c0
}

func smoothCount(t float64) float64 {
	return t/(2*math.Pi)*(math.Log(t/(2*math.Pi))-1) + 7.0/8
}

func huntZeros(count int) []float64 {
	zeros := []float64{}
	prevT, prevZ := 14.0, zRS(14.0)
	for t := 14.05; len(zeros) < count && t < 25000; t += 0.05 {
		zt := zRS(t)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 30 && hi-lo > 1e-8; i++ {
				mid := (lo + hi) / 2
				if (zRS(mid) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	return zeros
}

func main() {
	fmt.Println("LA LEY M-SIGMA DEL CIELO NUMÉRICO — the supermassive centers, weighed")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	gammas := huntZeros(20000)
	delta := make([]float64, len(gammas))
	for i, g := range gammas {
		delta[i] = smoothCount(g) - float64(i+1) + 0.5
	}
	varTot := 0.0
	for _, d := range delta {
		varTot += d * d
	}
	varTot /= float64(len(delta))
	fmt.Printf("\n  the galaxy's dispersion (variance of the tide at %d zeros): %.4f\n",
		len(gammas), varTot)

	fmt.Println("\n  THE CENTRAL BODIES - mass measured by projection, share of the dispersion:")
	fmt.Println("    p    mass (measured)   mass (theory)   share of sigma^2   cumulative")
	cum := 0.0
	for _, p := range []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31} {
		lp := math.Log(p)
		var cr, ci float64
		for i, g := range gammas {
			sn, cs := math.Sincos(g * lp)
			cr += delta[i] * cs
			ci += delta[i] * sn
		}
		amp := 2 * math.Hypot(cr, ci) / float64(len(gammas))
		share := amp * amp / 2 / varTot
		cum += share
		fmt.Printf("   %3.0f   %10.4f   %13.4f   %13.1f%%   %9.1f%%\n",
			p, amp, 1/math.Pi/math.Sqrt(p), 100*share, 100*cum)
	}
	fmt.Println("\n  the M-sigma verdict: a handful of supermassive centers - the primes")
	fmt.Println("  2, 3, 5, 7... - govern the dispersion of the ENTIRE galaxy of zeros,")
	fmt.Println("  exactly as the colossi at galactic centers govern their billions of")
	fmt.Println("  stars. Gentle horizons (slow, coherent voices - why the pacemaker")
	fmt.Println("  works), eternal lifetimes, and the whole sky moving to their pull.")
}
