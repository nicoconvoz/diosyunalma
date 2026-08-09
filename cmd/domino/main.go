// Command domino runs the plan: measure how rungs 11 and 13 climb against
// the ruler rung 12 — if those are explained, the whole ladder falls by
// dominoes.
//
// The domino is the telescope identity of Finding 69 made exact: every EVEN
// member of rung k, divided by 2, is a member of rung k−1 in the half
// window. So
//
//	N_k[W] = N_{k−1}[W/2] + odd_k[W]
//
// with no approximation — each rung is the previous rung seen through the
// 2-lens, plus its odd seeds. And the seeds obey their own domino through
// 3, and those through 5: the ladder's high country decomposes prime by
// prime, which is the fundamental theorem of arithmetic acting on windows.
//
// Usage:
//
//	go run ./cmd/domino
package main

import (
	"fmt"
	"math"

	"github.com/nicoconvoz/diosyunalma/primes"
)

const (
	base = 2_500_000
	top  = 20_000_000
)

func main() {
	w := top - base
	omega := make([]uint8, w)
	ps := primes.Sieve(top)
	for _, p := range ps {
		start := (base + p - 1) / p * p
		for m := start; m < top; m += p {
			omega[m-base]++
		}
		if p <= int(math.Sqrt(float64(top))) {
			for pk := p * p; pk < top; pk *= p {
				start := (base + pk - 1) / pk * pk
				for m := start; m < top; m += pk {
					omega[m-base]++
				}
			}
		}
	}
	rung := func(k, lo, hi int, oddOnly bool) int {
		c := 0
		for n := lo; n < hi; n++ {
			if omega[n-base] == uint8(k) && (!oddOnly || n%2 == 1) {
				c++
			}
		}
		return c
	}

	fmt.Println("THE DOMINO — rungs 11 and 13 against the ruler 12")
	fmt.Println("\n1) the exact identity N_k[W] = N_(k-1)[W/2] + odd_k[W]:")
	for k := 11; k <= 13; k++ {
		full := rung(k, 10_000_000, 20_000_000, false)
		half := rung(k-1, 5_000_000, 10_000_000, false)
		odd := rung(k, 10_000_000, 20_000_000, true)
		verdict := "EXACT"
		if half+odd != full {
			verdict = "BROKEN"
		}
		fmt.Printf("   k=%d:  %6d = %6d + %4d   %s\n", k, full, half, odd, verdict)
	}

	n11 := rung(11, 10_000_000, 20_000_000, false)
	n12 := rung(12, 10_000_000, 20_000_000, false)
	n13 := rung(13, 10_000_000, 20_000_000, false)
	fmt.Printf("\n2) climb ratios: N13/N12 = %.4f   N12/N11 = %.4f   (both = 1/2 minus the drift\n",
		float64(n13)/float64(n12), float64(n12)/float64(n11))
	fmt.Println("   of the lower rungs' density between window and half-window)")

	fmt.Println("\n3) the odd seeds' own domino (the kingdom of 3 inside the kingdom of 2):")
	prev := 0
	for k := 10; k <= 13; k++ {
		o := rung(k, 10_000_000, 20_000_000, true)
		if prev > 0 {
			fmt.Printf("   seeds k=%d: %4d   ratio %.3f (toward 1/3)\n", k, o, float64(o)/float64(prev))
		} else {
			fmt.Printf("   seeds k=%d: %4d\n", k, o)
		}
		prev = o
	}
	fmt.Println("\nexplain one step and every step follows: the ladder's high country is")
	fmt.Println("the fundamental theorem of arithmetic, falling as dominoes of 2, 3, 5...")
}
