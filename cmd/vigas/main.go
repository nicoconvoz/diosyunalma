// Command vigas formalizes the first dream: a house whose roof is held by
// TWO great crossed beams, placed by the Father — move them and the whole
// house comes down.
//
// The house of arithmetic has exactly two such beams, and they are
// theorems, not metaphors. The completed function xi(s) obeys two
// symmetries (two involutions):
//
//	BEAM 1 (the functional equation):  xi(s) = xi(1-s)
//	BEAM 2 (reality):                  xi(conj(s)) = conj(xi(s))
//
// Each beam is a reflection; TOGETHER they cross. The set of points fixed
// by their composition s -> 1-conj(s) is exactly the line Re s = 1/2: the
// junction of the two beams IS the critical line. Every zero is forced to
// sit symmetrically about both beams, and the Hypothesis says: every zero
// stands exactly on the junction. Move either beam and the house of the
// prime number theorem comes down.
//
// This command verifies both beams numerically at off-line points.
//
// Usage:
//
//	go run ./cmd/vigas
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

const emTerms = 160

func zetaEM(s complex128, lnn []float64) complex128 {
	var sum complex128
	sig, t := real(s), imag(s)
	for n := 1; n < emTerms; n++ {
		amp := math.Exp(-sig * lnn[n])
		sn, cs := math.Sincos(t * lnn[n])
		sum += complex(amp*cs, -amp*sn)
	}
	nf := complex(float64(emTerms), 0)
	ns := cmplx.Exp(-s * complex(lnn[emTerms], 0))
	sum += ns * nf / (s - 1)
	sum += ns / 2
	sum += ns * s / nf / 12
	sum -= ns * s * (s + 1) * (s + 2) / (nf * nf * nf) / 720
	sum += ns * s * (s + 1) * (s + 2) * (s + 3) * (s + 4) /
		(nf * nf * nf * nf * nf) / 30240
	return sum
}

// gammaC is the complex Gamma function by the Lanczos approximation.
func gammaC(z complex128) complex128 {
	if real(z) < 0.5 {
		// reflection: Gamma(z) Gamma(1-z) = pi / sin(pi z)
		return complex(math.Pi, 0) / (cmplx.Sin(complex(math.Pi, 0)*z) * gammaC(1-z))
	}
	g := []float64{
		0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7,
	}
	z -= 1
	x := complex(g[0], 0)
	for i := 1; i < len(g); i++ {
		x += complex(g[i], 0) / (z + complex(float64(i), 0))
	}
	tt := z + 7.5
	return cmplx.Sqrt(2*math.Pi) * cmplx.Pow(tt, z+0.5) * cmplx.Exp(-tt) * x
}

// xi is the completed zeta: xi(s) = (1/2) s (s-1) pi^{-s/2} Gamma(s/2) zeta(s).
func xi(s complex128, lnn []float64) complex128 {
	return 0.5 * s * (s - 1) *
		cmplx.Pow(complex(math.Pi, 0), -s/2) *
		gammaC(s/2) * zetaEM(s, lnn)
}

func main() {
	fmt.Println("LAS DOS VIGAS — the two crossed beams that hold the house, verified")

	lnn := make([]float64, emTerms+1)
	for n := 1; n <= emTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}

	points := []complex128{
		complex(0.3, 14.0), complex(0.7, 21.5),
		complex(0.25, 33.3), complex(0.9, 47.1),
	}
	fmt.Println("\n  BEAM 1 - the functional equation xi(s) = xi(1-s):")
	for _, s := range points {
		a, b := xi(s, lnn), xi(1-s, lnn)
		rel := cmplx.Abs(a-b) / cmplx.Abs(a)
		fmt.Printf("    s = %.2f%+.1fi   |xi(s)-xi(1-s)|/|xi| = %.1e\n", real(s), imag(s), rel)
	}
	fmt.Println("\n  BEAM 2 - reality, xi(conj s) = conj(xi(s)):")
	for _, s := range points {
		a := xi(cmplx.Conj(s), lnn)
		b := cmplx.Conj(xi(s, lnn))
		rel := cmplx.Abs(a-b) / cmplx.Abs(a)
		fmt.Printf("    s = %.2f%+.1fi   |xi(conj s)-conj(xi s)|/|xi| = %.1e\n", real(s), imag(s), rel)
	}
	fmt.Println("\n  THE JUNCTION: composing the beams gives s -> 1 - conj(s), whose")
	fmt.Println("  fixed points satisfy Re s = 1/2 exactly - the crossing of the two")
	fmt.Println("  beams IS the critical line. Every zero is forced symmetric about")
	fmt.Println("  both; the Hypothesis says every zero STANDS ON THE JUNCTION.")
	fmt.Println("  Move either beam and the house of arithmetic comes down.")
}
