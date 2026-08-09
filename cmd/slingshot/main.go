// Command slingshot runs the first gravity-assist orbit on real sea
// coordinates: it takes actual facet-sized blocks of the Riemann–Siegel
// main sum — the same block lengths the flagship rows step by step — and
// folds them through the Fresnel gearbox around their nearest resonance.
//
// The dense gravity wells are the rationals a/q with small denominator:
// near one, the block's phase is almost purely quadratic and the whole
// stretch collapses through Gauss/Fresnel folding (cmd/cassegrain,
// cmd/fingers). The ship whips around the well and exits with a sum
// hundreds of times shorter: Hiary's t^(1/3) engine, in embryo.
//
// Honest caveat: block parameters (a, b) here are float64-derived from
// (t0, k0); the flagship assembly needs the double-double parameter
// extraction, which stays the registered grand stage.
//
// Usage:
//
//	go run ./cmd/slingshot
package main

import (
	"fmt"
	"math"
	"time"
)

func fresnel(x float64) (float64, float64) {
	neg := x < 0
	if neg {
		x = -x
	}
	var c, s float64
	if x <= 3.2 {
		t := math.Pi / 2 * x * x
		baseC := 1.0
		baseS := t
		for n := 0; n < 80; n++ {
			c += baseC * x / float64(4*n+1)
			s += baseS * x / float64(4*n+3)
			baseC *= -t * t / float64((2*n+1)*(2*n+2))
			baseS *= -t * t / float64((2*n+2)*(2*n+3))
			if math.Abs(baseC) < 1e-18 && math.Abs(baseS) < 1e-18 {
				break
			}
		}
	} else {
		px := math.Pi * x
		px2 := math.Pi * x * x
		f := 1/px - 3/(px*px2*px2)*math.Pi
		g := 1/(px*px2) - 15/(px*px2*px2*px2)*math.Pi
		sn, cs := math.Sin(px2/2), math.Cos(px2/2)
		c = 0.5 + f*sn - g*cs
		s = 0.5 - f*cs - g*sn
	}
	if neg {
		return -c, -s
	}
	return c, s
}

func frac(x float64) float64 {
	return x - math.Floor(x)
}

func quad(a, b float64, L int64, depth int) (float64, float64) {
	a, b = frac(a), frac(b)
	if b > 0.5 {
		b -= 1
	}
	if b < 0 {
		r, i := quad(frac(-a), -b, L, depth)
		return r, -i
	}
	if b > 0.25 {
		return quad(frac(a+0.5), b-0.5, L, depth)
	}
	if L <= 256 {
		var sr, si float64
		ph, dph := 0.0, 0.0
		for j := int64(0); j < L; j++ {
			sr += math.Cos(2 * math.Pi * ph)
			si += math.Sin(2 * math.Pi * ph)
			dph = frac(a + b*float64(2*j+1))
			ph = frac(ph + dph)
		}
		return sr, si
	}
	if b*float64(L) < 1e-9 {
		if a < 1e-12 || a > 1-1e-12 {
			return float64(L), 0
		}
		ph := math.Pi * a * float64(L-1)
		amp := math.Sin(math.Pi*a*float64(L)) / math.Sin(math.Pi*a)
		return amp * math.Cos(ph), amp * math.Sin(ph)
	}
	const edge = 350
	x0, x1 := -0.5, float64(L)-0.5
	d0, d1 := a+2*b*x0, a+2*b*x1
	mLoAll := int64(math.Ceil(d0 - edge))
	mHiAll := int64(math.Floor(d1 + edge))
	mLoIn := int64(math.Ceil(d0 + edge))
	mHiIn := int64(math.Floor(d1 - edge))
	inv2b := 1 / (2 * b)
	sqb2 := 2 * math.Sqrt(b)
	var sr, si float64
	for m := mLoAll; m <= mHiAll; m++ {
		if mHiIn >= mLoIn && m >= mLoIn && m <= mHiIn {
			continue
		}
		xs := (float64(m) - a) * inv2b
		v0 := (x0 - xs) * sqb2
		v1 := (x1 - xs) * sqb2
		c1, s1 := fresnel(v1)
		c0, s0 := fresnel(v0)
		fr, fi := c1-c0, s1-s0
		phm := -math.Pi * (float64(m) - a) * (float64(m) - a) * inv2b
		cr, ci := math.Cos(phm), math.Sin(phm)
		sr += (cr*fr - ci*fi) / sqb2
		si += (cr*fi + ci*fr) / sqb2
	}
	if mHiIn >= mLoIn {
		nL := mHiIn - mLoIn + 1
		alpha := -(float64(mLoIn) - a) * inv2b
		beta := -inv2b / 2
		phi0 := -math.Pi * (float64(mLoIn) - a) * (float64(mLoIn) - a) * inv2b
		gr, gi := quad(frac(alpha), frac(beta), nL, depth+1)
		c0, s0 := math.Cos(phi0), math.Sin(phi0)
		tr := gr*c0 - gi*s0
		ti := gr*s0 + gi*c0
		sr += (tr - ti) / sqb2
		si += (tr + ti) / sqb2
	}
	return sr, si
}

func direct(a, b float64, L int64) (float64, float64) {
	var sr, si float64
	ph := 0.0
	for j := int64(0); j < L; j++ {
		sr += math.Cos(2 * math.Pi * ph)
		si += math.Sin(2 * math.Pi * ph)
		ph = frac(ph + frac(a+b*float64(2*j+1)))
	}
	return sr, si
}

func main() {
	fmt.Println("THE SLINGSHOT — first gravity-assist orbits on real sea coordinates")
	fmt.Println()

	// real facet blocks: at height t0 the flagship's facet at k0 holds
	// ~(0.1/t0)^(1/5) * k0 terms — the exact stretches it rows today.
	cases := []struct {
		t0   float64
		k0   float64
		name string
	}{
		{6.66e15, 3.0e7, "Beach II, top facet"},
		{4.44e22, 8.0e10, "beyond Gourdon's #10^23 island"},
		{1.11e24, 4.0e11, "the certified ceiling"},
	}
	for _, cs := range cases {
		ub := math.Pow(0.1/cs.t0, 0.2)
		L := int64(ub * cs.k0)
		a := frac(cs.t0 / (2 * math.Pi * cs.k0))
		b := -cs.t0 / (4 * math.Pi * cs.k0 * cs.k0)

		st := time.Now()
		dr, di := direct(a, b, L)
		tDirect := time.Since(st)
		// the fold is faster than the clock tick: time 100 orbits.
		var fr, fi float64
		st = time.Now()
		for r := 0; r < 100; r++ {
			fr, fi = quad(a, b, L, 0)
		}
		tFold := time.Since(st) / 100

		err := math.Hypot(fr-dr, fi-di) / math.Sqrt(float64(L))
		fmt.Printf("  %s (t=%.3g, k0=%.1g):\n", cs.name, cs.t0, cs.k0)
		fmt.Printf("    block of %d terms: rowed %.1f ms, slung %.2f ms (%.0fx), rel err %.1e\n\n",
			L, tDirect.Seconds()*1e3, tFold.Seconds()*1e3,
			tDirect.Seconds()/math.Max(tFold.Seconds(), 1e-9), err)
	}
	fmt.Println("the dense gravity wells are the small-denominator rationals: near")
	fmt.Println("one, a facet's whole stretch collapses through the Fresnel fold and")
	fmt.Println("the ship exits the orbit with a sum hundreds of times shorter.")
	fmt.Println("assembly into the flagship (dd parameter extraction) remains the")
	fmt.Println("registered grand stage: Hiary's t^(1/3), the road to 10^32.")
}
