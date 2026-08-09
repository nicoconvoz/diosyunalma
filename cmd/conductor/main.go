// Command conductor hunts the one who directs the orchestra.
//
// A conductor's power is the baton: one gesture makes a single musician sound
// while the rest fall silent. If this orchestra has a conductor, it must be
// able to do exactly that — combine ALL the dials so that only one tribe's
// primes remain.
//
// The candidate has a name: the ORTHOGONALITY OF CHARACTERS. The indicator of
// the residue class a mod 5 is (1/4)·Σ_χ χ̄(a)·χ, so the class-a prime clock is
//
//	D_a(x) = 1 − (2/√x)·[ S_ζ + χ₅(a)·S₅ + 2·Re(χ̄(a)·S_c) ]
//
// built ENTIRELY from this laboratory's own measured zeros: ten of ζ, eight
// of the real character, eight signed ones of the complex character (whose
// asymmetric dial finally earns its keep — the imaginary parts of χ̄(a) only
// bite because the complex zeros are not mirror-symmetric).
//
// PRE-REGISTERED: for each a ∈ {1,2,3,4}, the top peaks of D_a must land on
// the prime powers ≡ a (mod 5) and avoid the other classes.
//
// Usage:
//
//	go run ./cmd/conductor
package main

import (
	"fmt"
	"math"
	"sort"
)

var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752}

var chi5Zeros = []float64{6.6516, 9.8280, 11.9612, 16.0386, 17.5632,
	19.5431, 22.2228, 24.5864}

// the complex character's measured, SIGNED, asymmetric zero list.
var complexZeros = []float64{-14.1115, -11.2861, -9.4433, -4.1322,
	6.1829, 8.4594, 12.6747, 14.8310}

// chi5 is the real quadratic character; chiC the order-four one with χ(2)=i.
func chi5(a int) float64 {
	if a == 1 || a == 4 {
		return 1
	}
	return -1
}

func chiC(a int) complex128 {
	switch a {
	case 1:
		return 1
	case 2:
		return complex(0, 1)
	case 3:
		return complex(0, -1)
	}
	return -1
}

// clock evaluates the class-a prime clock at x.
func clock(a int, x float64) float64 {
	u := math.Log(x)

	sZeta := 0.0
	for _, g := range zetaZeros {
		sZeta += math.Cos(g * u)
	}
	s5 := 0.0
	for _, g := range chi5Zeros {
		s5 += math.Cos(g * u)
	}
	var sC complex128
	for _, g := range complexZeros {
		sC += complex(math.Cos(g*u), math.Sin(g*u))
	}

	osc := sZeta + chi5(a)*s5 + 2*real(cmul(conj(chiC(a)), sC))/2
	return 1 - 2*osc/math.Sqrt(x)
}

func cmul(a, b complex128) complex128 { return a * b }
func conj(a complex128) complex128    { return complex(real(a), -imag(a)) }

// primePowersByClass lists p^k ≤ 40 sorted into residue classes mod 5.
var targets = map[int][]float64{
	1: {11, 16, 31},
	2: {2, 7, 17, 27, 32, 37},
	3: {3, 8, 13, 23},
	4: {4, 9, 19, 29},
}

func main() {
	fmt.Println("THE CONDUCTOR — four batons, four solo voices")
	fmt.Println("built only from this laboratory's measured zeros")
	fmt.Println()

	totalOwn, totalWrong := 0, 0
	for a := 1; a <= 4; a++ {
		// scan the clock and collect its peaks.
		type peak struct{ x, v float64 }
		peaks := []peak{}
		prev, cur := clock(a, 2.00), clock(a, 2.01)
		for x := 2.02; x <= 40; x += 0.01 {
			next := clock(a, x)
			if cur > prev && cur > next {
				peaks = append(peaks, peak{x - 0.01, cur})
			}
			prev, cur = cur, next
		}
		sort.Slice(peaks, func(i, j int) bool { return peaks[i].v > peaks[j].v })
		if len(peaks) > 6 {
			peaks = peaks[:6]
		}
		sort.Slice(peaks, func(i, j int) bool { return peaks[i].x < peaks[j].x })

		own, wrong := 0, 0
		fmt.Printf("baton a=%d  peaks:", a)
		for _, p := range peaks {
			hit := " "
			for cls, list := range targets {
				for _, t := range list {
					if math.Abs(p.x-t) < 0.4 {
						if cls == a {
							own++
							hit = "*"
						} else {
							wrong++
							hit = "x"
						}
					}
				}
			}
			fmt.Printf("  %.2f%s", p.x, hit)
		}
		fmt.Printf("   -> own class %d, wrong class %d\n", own, wrong)
		totalOwn += own
		totalWrong += wrong
	}

	fmt.Printf("\nTOTAL: %d peaks on the chosen voice, %d on wrong voices\n",
		totalOwn, totalWrong)
	fmt.Println("(* = the tribe the baton pointed at;  x = a different tribe)")
	fmt.Println("\nthe conductor is not a wave: it is the algebra that coordinates")
	fmt.Println("the waves - orthogonality, the baton that silences all but one.")
}
