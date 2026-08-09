// Command lastimage builds the mirror in reverse.
//
// The flash: instead of walking the mirror from the beginning, go to the
// point where the image is lost, climb ONE step back to the last image
// still visible, and take the snapshot THERE. That rule exists: it is the
// optimal truncation of asymptotic series (Poincaré, Stokes, Berry) — a
// divergent edge-series' terms shrink, bottom out, and grow; stopping at
// the smallest term leaves an exponentially small error, and walking past
// it destroys everything.
//
// Two demonstrations:
//
//  1. THE PRINCIPLE — Stirling's series at x = 2 (truth: ln Γ(2) = 0):
//     the images shrink to step 7 and then grow; the total error bottoms
//     exactly at the last visible image.
//
//  2. THE APPLICATION — one extra snapshot step (the C₁ correction, built
//     by numerical differentiation of C₀) on the folded mirror at height
//     100,000: every zero improves ~150-fold, from seven to nine decimals.
//
// Usage:
//
//	go run ./cmd/lastimage
package main

import (
	"fmt"
	"math"
)

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 +
		1/(48*t) + 7/(5760*t*t*t)
}

func c0(p float64) float64 {
	return math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
}

func d3c0(p float64) float64 {
	const h = 0.008
	return (c0(p+2*h) - 2*c0(p+h) + 2*c0(p-h) - c0(p-2*h)) / (2 * h * h * h)
}

func z(t float64, deep bool) float64 {
	a := math.Sqrt(t / (2 * math.Pi))
	n := int(a)
	p := a - float64(n)
	th := theta(t)
	s := 0.0
	for k := 1; k <= n; k++ {
		s += math.Cos(th-t*math.Log(float64(k))) / math.Sqrt(float64(k))
	}
	s *= 2
	sign := 1.0
	if (n-1)%2 == 1 {
		sign = -1
	}
	tau := t / (2 * math.Pi)
	corr := c0(p)
	if deep {
		corr += -d3c0(p) / (96 * math.Pi * math.Pi) / math.Sqrt(tau)
	}
	return s + sign*math.Pow(tau, -0.25)*corr
}

func zeroNear(a, b float64, deep bool) float64 {
	lo, hi := a, b
	zlo := z(lo, deep)
	for i := 0; i < 60 && hi-lo > 2e-10; i++ {
		mid := (lo + hi) / 2
		zm := z(mid, deep)
		if (zlo < 0) != (zm < 0) {
			hi = mid
		} else {
			lo, zlo = mid, zm
		}
	}
	return (lo + hi) / 2
}

func main() {
	fmt.Println("THE LAST IMAGE — building the mirror in reverse")

	fmt.Println("\n1) the principle: Stirling at x = 2 (truth: ln Gamma(2) = 0)")
	fmt.Println("   step   image size    total error")
	bern := map[int]float64{2: 1.0 / 6, 4: -1.0 / 30, 6: 1.0 / 42, 8: -1.0 / 30,
		10: 5.0 / 66, 12: -691.0 / 2730, 14: 7.0 / 6, 16: -3617.0 / 510, 18: 43867.0 / 798}
	x := 2.0
	s := (x-0.5)*math.Log(x) - x + 0.5*math.Log(2*math.Pi)
	for n := 1; n <= 9; n++ {
		t := bern[2*n] / (float64(2*n) * float64(2*n-1) * math.Pow(x, float64(2*n-1)))
		s += t
		fmt.Printf("    %d     %.2e      %.2e\n", n, math.Abs(t), math.Abs(s))
	}
	fmt.Println("   the images shrink to step 7 and then GROW: stop at the last visible")
	fmt.Println("   one - the error bottoms exactly there. walk past it and all is lost.")

	fmt.Println("\n2) the application: one extra snapshot step on the mirror at 100,000")
	truth := []float64{100000.7437234872, 100001.1805583423,
		100002.0937225957, 100002.5172821591}
	brackets := [][2]float64{{100000.6, 100000.9}, {100001.0, 100001.3},
		{100001.9, 100002.2}, {100002.4, 100002.6}}
	fmt.Println("   true zero            error C0 only   error with C1")
	for i, tr := range truth {
		z0 := zeroNear(brackets[i][0], brackets[i][1], false)
		z1 := zeroNear(brackets[i][0], brackets[i][1], true)
		fmt.Printf("   %.10f   %.2e        %.2e\n", tr, math.Abs(z0-tr), math.Abs(z1-tr))
	}
	fmt.Println("\none step back from the vanishing point, snapshot there: the mirror")
	fmt.Println("gains two more decimals - the reverse construction, working.")
}
